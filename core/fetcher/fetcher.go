package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/config"
	"github.com/Tsisar/solana-indexer/core/parser"
	"github.com/Tsisar/solana-indexer/core/utils"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/Tsisar/solana-indexer/subgraph"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"sync"
)

var (
	runMu      sync.Mutex                        // Mutex to guard access to isRunning and hasPending flags
	isRunning  bool                              // Indicates whether the fetcher is currently running
	hasPending bool                              // Indicates if a retry is queued while fetcher is running
	client     = rpc.New(config.App.RPCEndpoint) // RPC client used for querying the Solana blockchain
)

// Start is the main entry point for the fetcher module.
// It begins a fetch cycle without resuming from the last signature.
func Start(ctx context.Context, db *storage.Gorm) {
	requestRun(ctx, db, "Start", false)
}

// Resume restarts the fetcher with resume=true,
// allowing it to continue from the last saved signature per program.
func Resume(ctx context.Context, db *storage.Gorm) {
	requestRun(ctx, db, "Resume", true)
}

// requestRun ensures only one fetcher instance runs at a time.
// If another run is already in progress, it optionally queues a retry.
func requestRun(ctx context.Context, db *storage.Gorm, label string, resume bool) {
	runMu.Lock()
	if isRunning {
		if hasPending {
			log.Infof("%s: fetcher already running and a retry is already queued — skipping", label)
			runMu.Unlock()
			return
		}
		hasPending = true
		log.Infof("%s: fetcher already running — queuing one retry", label)
		runMu.Unlock()
		return
	}
	isRunning = true
	runMu.Unlock()

	go runFetcher(ctx, db, label, resume)
}

// runFetcher handles a complete run of the fetcher,
// and reruns once if a retry was queued while processing.
func runFetcher(ctx context.Context, db *storage.Gorm, label string, resume bool) {
	for {
		log.Infof("%s: fetcher started", label)
		if err := fetch(ctx, db, resume); err != nil {
			subgraph.MapError(ctx, db, err)
			log.Fatalf("%s failed: %v", label, err)
		}

		runMu.Lock()
		if hasPending {
			hasPending = false
			runMu.Unlock()
			log.Infof("%s: rerunning due to queued retry", label)
			continue
		}
		isRunning = false
		runMu.Unlock()
		break
	}
}

// fetch orchestrates the entire data fetching and parsing process:
// 1. Fetch historical signatures,
// 2. Fetch full transaction JSONs,
// 3. Parse saved transactions.
func fetch(ctx context.Context, db *storage.Gorm, resume bool) error {
	if err := fetchHistoricalSignatures(ctx, db, resume); err != nil {
		return fmt.Errorf("failed to fetch historical signatures: %w", err)
	}
	if err := fetchRawTransactions(ctx, db); err != nil {
		return fmt.Errorf("failed to fetch full transactions: %w", err)
	}
	if err := parser.ParseSavedTransactions(ctx, db, resume); err != nil {
		return fmt.Errorf("failed to parse transactions: %w", err)
	}
	return nil
}

// fetchHistoricalSignatures retrieves transaction signatures for each program
// using paginated RPC requests. If resume is enabled, it will stop fetching
// once the last known signature is reached.
func fetchHistoricalSignatures(ctx context.Context, db *storage.Gorm, resume bool) error {
	programs := config.App.Programs
	r := config.App.EnableSignatureResume || resume

	for _, program := range programs {
		sigs, err := fetchHistoricalSignaturesForAddress(ctx, db, program, r)
		if err != nil {
			return fmt.Errorf("failed to fetch signatures for %s: %w", program, err)
		}

		for _, sig := range sigs {
			signatureStr := sig.Signature.String()
			transaction := core.Transaction{
				Signature: signatureStr,
				Slot:      sig.Slot,
				BlockTime: utils.BlockTime(sig.BlockTime),
			}

			if err := db.SaveTransaction(ctx, &transaction, program); err != nil {
				return fmt.Errorf("failed to save transaction %s: %w", signatureStr, err)
			}
		}
		log.Infof("Fetched %d signatures for program %s", len(sigs), program)
	}
	return nil
}

// fetchHistoricalSignaturesForAddress fetches all transaction signatures for a given address.
// It stops fetching once it reaches the last saved signature (if resume is enabled).
func fetchHistoricalSignaturesForAddress(ctx context.Context, db *storage.Gorm, address string, resume bool) ([]*rpc.TransactionSignature, error) {
	publicKey := solana.MustPublicKeyFromBase58(address)
	var before solana.Signature
	var until solana.Signature
	var result []*rpc.TransactionSignature

	if resume {
		if lastSigStr, err := db.GetLatestSavedSignature(ctx, address); err != nil {
			log.Errorf("get last saved signature failed: %v", err)
			return nil, err
		} else if lastSigStr != "" {
			until = solana.MustSignatureFromBase58(lastSigStr)
			log.Infof("Using last saved signature %s as lower bound for program %s", lastSigStr, address)
		}
	}

	for {
		log.Debugf("Fetching signatures from %s to %s", before, until)
		opts := &rpc.GetSignaturesForAddressOpts{
			Limit:      utils.Ptr(1000),
			Before:     before,
			Until:      until,
			Commitment: rpc.CommitmentConfirmed,
		}

		getSignaturesForAddressWithOpts := func() ([]*rpc.TransactionSignature, error) {
			return client.GetSignaturesForAddressWithOpts(ctx, publicKey, opts)
		}

		sigs, err := utils.Retry(getSignaturesForAddressWithOpts)
		if err != nil {
			return nil, fmt.Errorf("get signatures failed: %w", err)
		}

		if len(sigs) == 0 {
			break
		}

		result = append(result, sigs...)
		before = sigs[len(sigs)-1].Signature
	}

	return result, nil
}

// fetchRawTransactions retrieves full transaction data (JSON)
// for all transactions that do not yet have it (json_tx IS NULL).
func fetchRawTransactions(ctx context.Context, db *storage.Gorm) error {
	sigs, err := db.GetOrderedNoRawSignatures(ctx)
	if err != nil {
		return fmt.Errorf("get unparsed signatures failed: %w", err)
	}

	for _, sig := range sigs {
		txRes, err := FetchRawTransaction(ctx, sig)
		if err != nil {
			return fmt.Errorf("fetch raw transaction failed: %w", err)
		}

		raw, err := json.Marshal(txRes)
		if err != nil {
			return fmt.Errorf("marshal raw transaction failed: %w", err)
		}

		if err := db.UpdateTransactionRaw(ctx, sig, raw); err != nil {
			return fmt.Errorf("save transaction failed: %w", err)
		}
		log.Infof("Saved raw transaction: slot: %d tx: %s", txRes.Slot, sig)
	}
	return nil
}

// FetchRawTransaction retrieves the full transaction details from the RPC node
// for a given transaction signature.
func FetchRawTransaction(ctx context.Context, signature string) (*rpc.GetTransactionResult, error) {
	txSig := solana.MustSignatureFromBase58(signature)

	getTransactionResult := func() (*rpc.GetTransactionResult, error) {
		return client.GetTransaction(ctx, txSig, &rpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: utils.Ptr(uint64(0)),
		})
	}

	return utils.Retry(getTransactionResult)
}
