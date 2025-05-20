package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/config"
	"github.com/Tsisar/solana-indexer/core/utils"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

var client = rpc.New(config.App.RPCEndpoint) // RPC client used for querying the Solana blockchain

// Start orchestrates the entire data fetching and parsing process:
// 1. Fetch historical signatures,
// 2. Fetch full transaction JSONs,
// 3. Parse saved transactions.
func Start(ctx context.Context, db *storage.Gorm, resume bool, done chan struct{}) error {
	defer close(done) // ensure the signal is sent even on error

	if err := fetchHistoricalSignatures(ctx, db, resume); err != nil {
		return fmt.Errorf("[fetcher] failed to fetch historical signatures: %w", err)
	}
	if err := fetchRawTransactions(ctx, db); err != nil {
		return fmt.Errorf("[fetcher] failed to fetch full transactions: %w", err)
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
			return fmt.Errorf("[fetcher] failed to fetch signatures for %s: %w", program, err)
		}

		for _, sig := range sigs {
			signatureStr := sig.Signature.String()
			transaction := core.Transaction{
				Signature: signatureStr,
				Slot:      sig.Slot,
				BlockTime: utils.BlockTime(sig.BlockTime),
			}

			if err := db.SaveTransaction(ctx, &transaction, program); err != nil {
				return fmt.Errorf("[fetcher] failed to save transaction %s: %w", signatureStr, err)
			}
		}
		log.Infof("[fetcher] Fetched %d signatures for program %s", len(sigs), program)
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
			log.Errorf("[fetcher] get last saved signature failed: %v", err)
			return nil, err
		} else if lastSigStr != "" {
			until = solana.MustSignatureFromBase58(lastSigStr)
			log.Infof("[fetcher] Using last saved signature %s as lower bound for program %s", lastSigStr, address)
		}
	}

	for {
		log.Debugf("[fetcher] Fetching signatures from %s to %s", before, until)
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
			return nil, fmt.Errorf("[fetcher] get signatures failed: %w", err)
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
		return fmt.Errorf("[fetcher] get unparsed signatures failed: %w", err)
	}

	for _, sig := range sigs {
		txRes, err := FetchRawTransaction(ctx, sig)
		if err != nil {
			return fmt.Errorf("[fetcher] fetch raw transaction failed: %w", err)
		}

		raw, err := json.Marshal(txRes)
		if err != nil {
			return fmt.Errorf("[fetcher] marshal raw transaction failed: %w", err)
		}

		if err := db.UpdateTransactionRaw(ctx, sig, raw); err != nil {
			return fmt.Errorf("[fetcher] save transaction failed: %w", err)
		}
		log.Infof("[fetcher] Saved raw transaction: slot: %d tx: %s", txRes.Slot, sig)
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
