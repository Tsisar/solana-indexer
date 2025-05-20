package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/config"
	"github.com/Tsisar/solana-indexer/core/fetcher"
	"github.com/Tsisar/solana-indexer/core/utils"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"gorm.io/datatypes"
)

// fetchTask represents a signature received from a WebSocket subscription
// along with the associated program it was observed in.
type fetchTask struct {
	Signature string
	Program   string
}

// Start initializes WebSocket subscriptions and starts processing
// of transactions fetched via WebSocket events.
// It ensures all programs are subscribed before signaling readiness via wsReady.
func Start(ctx context.Context, db *storage.Gorm, wsReady chan<- struct{}, realtimeStream chan<- string, errorChan chan<- error) error {
	// Connect to Solana WebSocket endpoint
	wsClient, err := ws.Connect(ctx, config.App.RPCWSEndpoint)
	if err != nil {
		return fmt.Errorf("[listener] failed to connect to WebSocket: %w", err)
	}

	fetchQueue := make(chan fetchTask, 1000)
	programCount := len(config.App.Programs)
	connected := make(chan struct{}, programCount)

	// Start transaction fetcher that processes incoming WebSocket events
	go func() {
		if err := fetchFromQueue(ctx, db, fetchQueue, realtimeStream); err != nil {
			errorChan <- fmt.Errorf("[listener] fetcher error: %v", err)
		}
	}()

	// Start WebSocket subscriptions for all configured programs
	for _, program := range config.App.Programs {
		publicKey := solana.MustPublicKeyFromBase58(program)

		go func(pid solana.PublicKey) {
			if err := watch(ctx, wsClient, pid, connected, fetchQueue); err != nil {
				errorChan <- fmt.Errorf("[listener] watch failed for %s: %w", pid.String(), err)
			}
		}(publicKey)
	}

	// Wait until all programs are subscribed or an error occurs
	readyCount := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-connected:
			readyCount++
			if readyCount == programCount {
				log.Info("[listener] All WebSocket subscriptions are active.")
				wsReady <- struct{}{}
				return nil
			}
		}
	}
}

// watch subscribes to transaction logs for a single Solana program
// and forwards the observed signatures to the fetchQueue.
func watch(ctx context.Context, wsClient *ws.Client, publicKey solana.PublicKey, readySignal chan<- struct{}, fetchQueue chan<- fetchTask) error {
	log.Infof("[listener] subscribing for %s...", publicKey.String())

	sub, err := wsClient.LogsSubscribeMentions(publicKey, rpc.CommitmentConfirmed)
	if err != nil {
		return fmt.Errorf("[listener] logs subscribe failed for %s: %w", publicKey.String(), err)
	}

	log.Infof("[listener] subscribed to program %s", publicKey.String())
	readySignal <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			sub.Unsubscribe()
			return nil
		default:
			msg, err := sub.Recv(ctx)
			if err != nil {
				log.Errorf("[listener] WebSocket recv failed for %s: %v", publicKey.String(), err)
				sub.Unsubscribe()
				return fmt.Errorf("[listener] recv failed for %s: %w", publicKey.String(), err)
			}

			if msg == nil {
				continue
			}

			// Push received signature into fetch queue for processing
			fetchQueue <- fetchTask{
				Signature: msg.Value.Signature.String(),
				Program:   publicKey.String(),
			}
		}
	}
}

// fetchFromQueue processes transactions sequentially from the fetchQueue:
// it ensures each transaction is fetched from RPC and stored in the DB,
// and then its signature is passed to the parser stream.
func fetchFromQueue(ctx context.Context, db *storage.Gorm, queue <-chan fetchTask, stream chan<- string) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case task := <-queue:
			if err := fetch(ctx, db, task.Program, task.Signature, stream); err != nil {
				return fmt.Errorf("[listener] fetch failed for transaction %s: %w", task.Signature, err)
			}
			log.Infof("[listener] Fetched transaction %s for program %s", task.Signature, task.Program)
		}
	}
}

// fetch loads a transaction from the Solana RPC,
// stores it in the database, and pushes the signature to the parsing stream.
func fetch(ctx context.Context, db *storage.Gorm, program, signature string, stream chan<- string) error {
	// Skip if already fetched
	fetched, err := db.IsRawFetched(ctx, signature)
	if err != nil {
		return fmt.Errorf("[listener] failed to check if transaction %s is fetched: %w", signature, err)
	}
	if fetched {
		log.Warnf("[listener] Transaction %s already fetched, skipping...", signature)
		err := db.AssociateTransactionWithProgram(ctx, signature, program)
		if err != nil {
			return fmt.Errorf("[listener] failed to associate transaction %s with program %s: %w", signature, program, err)
		}
		return nil
	}

	// Fetch raw transaction from RPC
	txRes, err := fetcher.FetchRawTransaction(ctx, signature)
	if err != nil {
		return fmt.Errorf("[listener] failed to fetch raw transaction %s: %w", signature, err)
	}

	// Serialize raw transaction to JSON
	raw, err := json.Marshal(txRes)
	if err != nil {
		return fmt.Errorf("[listener] failed to marshal raw transaction %s: %w", signature, err)
	}

	// Save to database
	transaction := core.Transaction{
		Signature: signature,
		Slot:      txRes.Slot,
		BlockTime: utils.BlockTime(txRes.BlockTime),
		JsonTx:    datatypes.JSON(raw),
	}

	if err := db.SaveTransaction(ctx, &transaction, program); err != nil {
		return fmt.Errorf("[listener] failed to save transaction %s: %w", signature, err)
	}

	// Push the signature to the realtime parsing stream
	select {
	case stream <- signature:
		return nil
	case <-ctx.Done():
		return nil
	}
}
