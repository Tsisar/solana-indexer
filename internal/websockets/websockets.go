package websockets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/config"
	"github.com/Tsisar/solana-indexer/internal/fetcher"
	"github.com/Tsisar/solana-indexer/internal/parser"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/Tsisar/solana-indexer/internal/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"time"
)

var (
	client = rpc.New(config.App.RPCEndpoint)
)

func Start(ctx context.Context, db *storage.Gorm) {
	wsClient, err := ws.Connect(ctx, config.App.RPCWSEndpoint)
	if err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}

	for _, programs := range config.App.Programs {
		publicKey := solana.MustPublicKeyFromBase58(programs)

		go func(pid solana.PublicKey, pidStr string) {
			if err := watch(ctx, db, wsClient, pid, pidStr); err != nil {
				log.Fatalf("Watch failed for %s: %v", pidStr, err)
			}
		}(publicKey, programs)
	}
}

func watch(ctx context.Context, db *storage.Gorm, wsClient *ws.Client, publicKey solana.PublicKey, programs string) error {
	log.Infof("WebSocket: watching program %s", programs)

	const maxFailures = 5 // TODO: make this configurable
	failures := 0

	for {
		sub, err := wsClient.LogsSubscribeMentions(publicKey, rpc.CommitmentConfirmed)
		if err != nil {
			failures++
			log.Errorf("logs subscribe failed for %s (attempt %d/%d): %v", programs, failures, maxFailures, err)

			if failures >= maxFailures {
				return fmt.Errorf("max reconnect attempts reached for %s", programs)
			}

			time.Sleep(10 * time.Second)
			continue
		}

		log.Infof("WebSocket: subscribed to program %s", programs)
		if failures > 0 {
			log.Infof("WebSocket: reconnected to program %s after %d failures", programs, failures)
			go fetcher.Resume(ctx, db)
			failures = 0
		}

		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				return nil

			default:
				msg, err := sub.Recv(ctx)
				if err != nil {
					log.Warnf("WebSocket recv failed for %s: %v", programs, err)
					sub.Unsubscribe()
					time.Sleep(5 * time.Second)
					break
				}

				if msg == nil {
					continue
				}

				log.Debugf("WebSocket: received message for %s: %s", programs, msg.Value.Signature)
				//go fetcher.Resume(ctx, db)
				go fetch(ctx, db, programs, msg.Value.Signature.String(), msg.Context.Slot)
			}
		}
	}
}

func fetch(ctx context.Context, db *storage.Gorm, program, signature string, slot uint64) {
	ok, err := db.IsParsed(ctx, signature)
	if err != nil {
		log.Errorf("Failed to check if transaction %s is parsed: %v", signature, err)
		return
	}
	if !ok {
		log.Warnf("Transaction %s already parsed, skipping", signature)
		return
	}

	if err := saveTransactions(ctx, db, program, signature, slot); err != nil {
		log.Errorf("Failed to save transaction %s: %v", signature, err)
		return
	}

	rawTx, err := fetchFullTransactions(ctx, db, signature)
	if err != nil {
		log.Errorf("Failed to fetch full transaction %s: %v", signature, err)
		return
	}

	if err := parser.ParseTransaction(ctx, db, rawTx, signature); err != nil {
		log.Errorf("failed to process tx %s: %v", signature, err)
		return
	}

	if err := db.MarkParsed(ctx, signature); err != nil {
		log.Errorf("failed to mark parsed %s: %v", signature, err)
	}
}

func saveTransactions(ctx context.Context, db *storage.Gorm, program, signature string, slot uint64) error {
	if err := db.SaveTransaction(ctx, signature, slot); err != nil {
		return fmt.Errorf("failed to save transaction %s: %w", signature, err)
	}

	if err := db.AssociateTransactionWithProgram(ctx, signature, program); err != nil {
		return fmt.Errorf("failed to associate program %s with transaction %s: %w", program, signature, err)
	}

	return nil
}

func fetchFullTransactions(ctx context.Context, db *storage.Gorm, sig string) ([]byte, error) {
	txSig := solana.MustSignatureFromBase58(sig)

	getTransactionResult := func() (*rpc.GetTransactionResult, error) {
		return client.GetTransaction(ctx, txSig, &rpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: utils.Ptr(uint64(0)),
		})
	}

	txRes, err := utils.Retry(getTransactionResult)
	if err != nil {
		return nil, fmt.Errorf("get transaction failed: %w", err)
	}

	raw, err := json.Marshal(txRes)
	if err != nil {
		return nil, fmt.Errorf("marshal tx failed: %w", err)
	}

	if err := db.UpdateTransactionRaw(ctx, sig, raw); err != nil {
		return nil, fmt.Errorf("save transaction failed: %w", err)
	}
	log.Infof("Saved raw transaction slot: %d, tx: %s", txRes.Slot, sig)

	return raw, nil
}
