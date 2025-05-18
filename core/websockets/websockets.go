package websockets

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/config"
	"github.com/Tsisar/solana-indexer/core/fetcher"
	"github.com/Tsisar/solana-indexer/core/parser"
	"github.com/Tsisar/solana-indexer/core/utils"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/Tsisar/solana-indexer/subgraph"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/ws"
	"gorm.io/datatypes"
	"time"
)

type fetchTask struct {
	Signature string
	Program   string
}

// Start initializes the WebSocket listeners and fetch queue processor.
// It establishes subscriptions to all configured Solana programs and starts
// a worker to process incoming transactions.
func Start(ctx context.Context, db *storage.Gorm) error {
	// Queue for handling incoming transactions via WebSocket
	fetchQueue := make(chan fetchTask, 1000)

	// Worker that waits until the DB is ready and then processes queued transactions.
	go waitAndProcess(ctx, db, fetchQueue)

	// Connect to Solana WebSocket RPC endpoint
	wsClient, err := ws.Connect(ctx, config.App.RPCWSEndpoint)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	// Subscribe to logs for each configured program
	for _, program := range config.App.Programs {
		publicKey := solana.MustPublicKeyFromBase58(program)
		go func(pid solana.PublicKey, pidStr string) {
			if err := watch(ctx, db, wsClient, pid, pidStr, fetchQueue); err != nil {
				subgraph.MapError(ctx, db, err)
				log.Fatalf("Watch failed for %s: %v", pidStr, err)
			}
		}(publicKey, program)
	}
	return nil
}

// watch subscribes to log messages for the given program and sends new transactions into the fetch queue.
// It handles reconnection and retry logic on failures.
func watch(ctx context.Context, db *storage.Gorm, wsClient *ws.Client, publicKey solana.PublicKey, program string, fetchQueue chan<- fetchTask) error {
	log.Infof("[watch:%s] watching...", program)

	const maxFailures = 10
	failures := 0

OUTER:
	for {
		// Subscribe to logs mentioning the given program
		sub, err := wsClient.LogsSubscribeMentions(publicKey, rpc.CommitmentConfirmed)
		if err != nil {
			failures++
			log.Errorf("logs subscribe failed for %s (attempt %d/%d): %v", program, failures, maxFailures, err)

			if failures >= maxFailures {
				return fmt.Errorf("max reconnect attempts reached for %s", program)
			}

			time.Sleep(10 * time.Second)
			continue OUTER
		}

		log.Infof("WebSocket: subscribed to program %s", program)
		if failures > 0 {
			log.Infof("WebSocket: reconnected to program %s after %d failures", program, failures)
			failures = 0
			// Trigger fetcher retry after reconnection
			fetcher.Resume(ctx, db)
		}

		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				return nil
			default:
				msg, err := sub.Recv(ctx)
				if err != nil {
					failures++
					log.Warnf("WebSocket recv failed for %s: %v", program, err)
					sub.Unsubscribe()
					time.Sleep(5 * time.Second)
					continue OUTER
				}

				if msg == nil {
					continue
				}

				// Send the signature to the processing queue
				fetchQueue <- fetchTask{
					Signature: msg.Value.Signature.String(),
					Program:   program,
				}
			}
		}
	}
}

// waitAndProcess waits for the database to become healthy and ready,
// then starts processing the transaction queue as items arrive.
func waitAndProcess(ctx context.Context, db *storage.Gorm, queue <-chan fetchTask) {
	log.Info("Waiting for database readiness before starting fetcher...")

	// Poll for DB readiness and health status
	for {
		status, reason, err := db.GetHealth(ctx)
		if err != nil {
			log.Errorf("Failed to get database health: %v", err)
		} else if status != "healthy" {
			log.Warnf("Database is not healthy: %s", reason)
			break
		}
		ready, err := db.IsReady(ctx)
		if err != nil {
			log.Errorf("IsReady check failed: %v", err)
		} else if ready {
			log.Info("Database is ready. Starting fetch processing.")
			break
		}
		time.Sleep(5 * time.Second)
	}

	// Process fetch tasks from the queue
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-queue:
			log.Debugf("Processing from queue: %s", task.Signature)
			fetch(ctx, db, task.Program, task.Signature)
		}
	}
}

// fetch performs the full lifecycle for handling a transaction:
// 1. Skips if already parsed
// 2. Fetches the raw transaction from RPC
// 3. Saves it in the database
// 4. Parses and stores events
// 5. Marks the transaction as parsed
func fetch(ctx context.Context, db *storage.Gorm, program, signature string) {
	parsed, err := db.IsParsed(ctx, signature)
	if err != nil {
		log.Errorf("Failed to check if transaction %s is parsed: %v", signature, err)
		return
	}
	if parsed {
		log.Warnf("Transaction %s already parsed, skipping", signature)
		return
	}

	// Fetch the transaction from RPC
	txRes, err := fetcher.FetchRawTransaction(ctx, signature)
	if err != nil {
		log.Errorf("Failed to fetch raw transaction %s: %v", signature, err)
		return
	}

	raw, err := json.Marshal(txRes)
	if err != nil {
		log.Errorf("Failed to marshal raw transaction %s: %v", signature, err)
		return
	}

	// Save the transaction into the DB
	transaction := core.Transaction{
		Signature: signature,
		Slot:      txRes.Slot,
		BlockTime: utils.BlockTime(txRes.BlockTime),
		JsonTx:    datatypes.JSON(raw),
	}

	if err := db.SaveTransaction(ctx, &transaction, program); err != nil {
		log.Errorf("Failed to save transaction %s: %v", signature, err)
		return
	}

	// Parse the transaction (logs + instructions)
	if err := parser.ParseTransaction(ctx, db, raw, signature); err != nil {
		log.Errorf("Failed to parse transaction %s: %v", signature, err)
		return
	}

	// Mark as parsed
	if err := db.MarkParsed(ctx, signature); err != nil {
		log.Errorf("failed to mark parsed %s: %v", signature, err)
	}
}
