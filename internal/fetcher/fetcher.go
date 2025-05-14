package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/config"
	"github.com/Tsisar/solana-indexer/internal/parser"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/Tsisar/solana-indexer/internal/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"sync"
)

var (
	runMu      sync.Mutex
	isRunning  bool
	hasPending bool
	client     = rpc.New(config.App.RPCEndpoint)
)

// Start is the main entry point for the fetcher module.
func Start(ctx context.Context, db *storage.Gorm) {
	requestRun(ctx, db, "Start", false)
}

// Resume restarts the fetcher with resume=true.
func Resume(ctx context.Context, db *storage.Gorm) {
	requestRun(ctx, db, "Resume", true)
}

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

func runFetcher(ctx context.Context, db *storage.Gorm, label string, resume bool) {
	for {
		log.Infof("%s: fetcher started", label)
		if err := fetch(ctx, db, resume); err != nil {
			log.Errorf("%s failed: %v", label, err)
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

func fetch(ctx context.Context, db *storage.Gorm, resume bool) error {
	if err := fetchHistoricalSignatures(ctx, db, resume); err != nil {
		return fmt.Errorf("failed to fetch historical signatures: %w", err)
	}
	if err := fetchFullTransactions(ctx, db); err != nil {
		return fmt.Errorf("failed to fetch full transactions: %w", err)
	}
	parser.Start(ctx, db)

	return nil
}

// fetchHistoricalSignatures retrieves signatures for a given program using paginated RPC calls.
// If signature resumption is enabled, it stops at the last saved signature (exclusive).
func fetchHistoricalSignatures(ctx context.Context, db *storage.Gorm, resume bool) error {
	programs := config.App.Programs
	r := config.App.EnableSignatureResume || resume

	// Collect all signatures for each program
	for _, program := range programs {
		sigs, err := fetchHistoricalSignaturesForAddress(ctx, db, program, r)
		if err != nil {
			return fmt.Errorf("failed to fetch signatures for %s: %w", program, err)
		}

		if err := db.SaveProgram(ctx, program); err != nil {
			return fmt.Errorf("failed to save program %s: %w", program, err)
		}

		for _, sig := range sigs {
			sigStr := sig.Signature.String()

			if err := db.SaveTransaction(ctx, sigStr, sig.Slot); err != nil {
				return fmt.Errorf("failed to save transaction %s: %w", sigStr, err)
			}
			if err := db.AssociateTransactionWithProgram(ctx, sigStr, program); err != nil {
				return fmt.Errorf("failed to associate program %s with transaction %s: %w", program, sigStr, err)
			}
		}
		log.Infof("Fetched %d signatures for program %s", len(sigs), program)
	}
	return nil
}

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

	// Pagination loop
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

// fetchFullTransactions retrieves full transaction data for all transactions
// that haven't yet been enriched with raw data (json_tx IS NULL).
// It marshals the RPC response and updates the existing transaction entry.
func fetchFullTransactions(ctx context.Context, db *storage.Gorm) error {
	// Get signatures for fetching
	sigs, err := db.GetOrderedNoRawSignatures(ctx)
	if err != nil {
		return fmt.Errorf("get unparsed signatures failed: %w", err)
	}

	for _, sig := range sigs {
		txSig := solana.MustSignatureFromBase58(sig)
		// Get full transaction data

		getTransactionResult := func() (*rpc.GetTransactionResult, error) {
			return client.GetTransaction(ctx, txSig, &rpc.GetTransactionOpts{
				Encoding:                       solana.EncodingBase64,
				Commitment:                     rpc.CommitmentConfirmed,
				MaxSupportedTransactionVersion: utils.Ptr(uint64(0)),
			})
		}

		txRes, err := utils.Retry(getTransactionResult)
		if err != nil {
			return fmt.Errorf("get transaction failed: %w", err)
		}

		raw, err := json.Marshal(txRes)
		if err != nil {
			return fmt.Errorf("marshal tx failed: %w", err)
		}

		if err := db.UpdateTransactionRaw(ctx, sig, raw); err != nil {
			return fmt.Errorf("save transaction failed: %w", err)
		}
		log.Infof("Saved raw transaction slot: %d, tx: %s", txRes.Slot, sig)
	}
	return nil
}
