package reader

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/config"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/Tsisar/solana-indexer/internal/utils"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var (
	client = rpc.New(config.App.RPCEndpoint)
)

// Start is the main entry point for the fetcher module.
func Start(ctx context.Context, db *storage.Gorm, signatureChannel chan string) {
	if err := fetchHistoricalSignatures(ctx, db, false); err != nil {
		log.Errorf("Failed to fetch historical signatures: %v", err)
	}
	if err := getSavedSignatures(ctx, db, signatureChannel); err != nil {
		log.Errorf("Failed to get saved signatures: %v", err)
	}
}

// Resume restarts the fetcher with resume=true.
func Resume(ctx context.Context, db *storage.Gorm, signatureChannel chan string) {
	if err := fetchHistoricalSignatures(ctx, db, true); err != nil {
		log.Errorf("Failed to fetch historical signatures: %v", err)
	}
	if err := getSavedSignatures(ctx, db, signatureChannel); err != nil {
		log.Errorf("Failed to get saved signatures: %v", err)
	}
}

func Listen(ctx context.Context, db *storage.Gorm, signatureChannel chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case signature := <-signatureChannel:
			slot, err := fetchTransaction(ctx, db, signature)
			if err != nil {
				log.Fatalf("Failed to get transaction: %v", err)
			}
			log.Infof("Saved raw transaction slot: %d, tx: %s", slot, signature)
		}
	}
}

func fetchTransaction(ctx context.Context, db *storage.Gorm, signature string) (uint64, error) {
	txSig := solana.MustSignatureFromBase58(signature)
	// Fetch full transaction data

	getTransactionResult := func() (*rpc.GetTransactionResult, error) {
		return client.GetTransaction(ctx, txSig, &rpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: utils.Ptr(uint64(0)),
		})
	}

	txRes, err := utils.Retry(getTransactionResult)
	if err != nil {
		return 0, fmt.Errorf("get transaction failed: %w", err)
	}

	raw, err := json.Marshal(txRes)
	if err != nil {
		return 0, fmt.Errorf("marshal tx failed: %w", err)
	}

	if err := db.UpdateTransactionRaw(ctx, signature, raw); err != nil {
		return 0, fmt.Errorf("save transaction failed: %w", err)
	}

	return txRes.Slot, nil
}

func getSavedSignatures(ctx context.Context, db *storage.Gorm, signatureChannel chan string) error {
	sigs, err := db.GetOrderedNoRawSignatures(ctx)
	if err != nil {
		return fmt.Errorf("get no raw signatures failed: %w", err)
	}
	for _, sig := range sigs {
		signatureChannel <- sig
	}
	return nil
}

// fetchHistoricalSignatures retrieves signatures for a given program using paginated RPC calls.
// If signature resumption is enabled, it stops at the last saved signature (exclusive).
func fetchHistoricalSignatures(ctx context.Context, db *storage.Gorm, resume bool) error {
	programs := config.App.Programs
	r := config.App.EnableSignatureResume || resume

	// Collect all signatures for each program
	for _, program := range programs {
		signatures, err := fetchHistoricalSignaturesForAddress(ctx, db, program, r)
		if err != nil {
			return fmt.Errorf("failed to fetch signatures for %s: %w", program, err)
		}
		if err := db.SaveProgram(ctx, program); err != nil {
			return fmt.Errorf("failed to save program %s: %w", program, err)
		}

		for _, signature := range signatures {
			signatureStr := signature.Signature.String()
			blockTime := int64(*signature.BlockTime)

			if err := db.SaveTransaction(ctx, signatureStr, signature.Slot, blockTime); err != nil {
				return fmt.Errorf("failed to save transaction %s: %w", signatureStr, err)
			}
		}
		log.Infof("Fetched %d signatures for program %s", len(signatures), program)
	}
	return nil
}

func fetchHistoricalSignaturesForAddress(ctx context.Context, db *storage.Gorm, address string, resume bool) ([]*rpc.TransactionSignature, error) {
	publicKey := solana.MustPublicKeyFromBase58(address)
	var before solana.Signature
	var until solana.Signature
	var result []*rpc.TransactionSignature

	if resume {
		if latestSignatureStr, err := db.GetLatestSavedSignature(ctx, address); err != nil {
			log.Errorf("get last saved signature failed: %v", err)
			return nil, err
		} else if latestSignatureStr != "" {
			until = solana.MustSignatureFromBase58(latestSignatureStr)
			log.Infof("Using last saved signature %s as lower bound for program %s", latestSignatureStr, address)
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
