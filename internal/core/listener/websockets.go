package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/config"
	"github.com/Tsisar/solana-indexer/internal/core/fetcher"
	"github.com/Tsisar/solana-indexer/internal/monitoring"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/Tsisar/solana-indexer/internal/storage/model/core"
	"github.com/Tsisar/solana-indexer/internal/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"gorm.io/datatypes"
	"sort"
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
func Start(ctx context.Context, db *storage.Gorm, receiveHandler func(signature string), readyHandler func(), errorHandler func(err error)) error {
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
		if err := fetchFromQueue(ctx, db, fetchQueue, receiveHandler); err != nil {
			errorHandler(fmt.Errorf("[listener] fetcher failed: %v", err))
		}
	}()

	// Start WebSocket subscriptions for all configured programs
	for _, program := range config.App.Programs {
		go func() {
			if err := watch(ctx, db, wsClient, program, connected, fetchQueue); err != nil {
				errorHandler(fmt.Errorf("[listener] watch failed for %s: %v", program, err))
			}
		}()
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
				readyHandler()
				return nil
			}
		}
	}
}

// watch subscribes to transaction logs for a single Solana program
// and forwards the observed signatures to the fetchQueue.
func watch(ctx context.Context, db *storage.Gorm, wsClient *ws.Client, program string, readySignal chan<- struct{}, fetchQueue chan<- fetchTask) error {
	log.Infof("[listener] subscribing for %s...", program)
	publicKey := solana.MustPublicKeyFromBase58(program)

	sub, err := wsClient.LogsSubscribeMentions(publicKey, rpc.CommitmentConfirmed)
	if err != nil {
		return fmt.Errorf("[listener] logs subscribe failed for %s: %w", program, err)
	}

	log.Infof("[listener] subscribed to program %s", program)

	// Check if we are not skipped transactions
	if err := checkSkipped(ctx, db, publicKey, fetchQueue); err != nil {
		return fmt.Errorf("[listener] check skipped failed for %s: %w", program, err)
	}
	readySignal <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			sub.Unsubscribe()
			return nil
		default:
			msg, err := sub.Recv(ctx)
			if err != nil {
				sub.Unsubscribe()
				return fmt.Errorf("[listener] recv failed for %s: %w", program, err)
			}

			if msg == nil {
				continue
			}

			// Push received signature into fetch queue for processing
			fetchQueue <- fetchTask{
				Signature: msg.Value.Signature.String(),
				Program:   program,
			}
		}
	}
}

// fetchFromQueue processes transactions sequentially from the fetchQueue:
// it ensures each transaction is fetched from RPC and stored in the DB,
// and then its signature is passed to the parser stream.
func fetchFromQueue(ctx context.Context, db *storage.Gorm, queue <-chan fetchTask, receiveHandler func(signature string)) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case task := <-queue:
			if err := fetch(ctx, db, task.Program, task.Signature, receiveHandler); err != nil {
				return fmt.Errorf("[listener] fetch failed for transaction %s: %w", task.Signature, err)
			}
		}
	}
}

// fetch loads a transaction from the Solana RPC,
// stores it in the database, and pushes the signature to the parsing stream.
func fetch(ctx context.Context, db *storage.Gorm, program, signature string, receiveHandler func(signature string)) error {
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
	log.Infof("[listener] Fetched transaction %s", signature)
	monitoring.ListenerCurrentSlot.Set(float64(txRes.Slot))
	receiveHandler(signature)
	return nil
}

func checkSkipped(ctx context.Context, db *storage.Gorm, publicKey solana.PublicKey, fetchQueue chan<- fetchTask) error {
	var before solana.Signature
	until, err := db.GetLatestSavedSignature(ctx, publicKey.String())
	if err != nil {
		return fmt.Errorf("[listener] get last saved signature failed: %v", err)
	}

	sigs, err := fetcher.GetSignaturesForRange(ctx, publicKey, before, until)
	if err != nil {
		return fmt.Errorf("[listener] get signatures failed: %v", err)
	}
	log.Infof("[listener] Received %d signatures for program %s", len(sigs), publicKey.String())

	if len(sigs) == 0 {
		log.Infof("[listener] No signatures found for program %s", publicKey.String())
		return nil
	}

	//Sort signatures by slot
	sort.Slice(sigs, func(i, j int) bool {
		return sigs[i].Slot < sigs[j].Slot
	})

	for _, sig := range sigs {
		log.Debugf("[listener] Block %d, adding skipped signature %s to fetch queue", sig.Slot, sig.Signature.String())

		fetchQueue <- fetchTask{
			Signature: sig.Signature.String(),
			Program:   publicKey.String(),
		}
	}
	return nil
}
