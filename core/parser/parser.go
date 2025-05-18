package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/gagliardetto/solana-go/rpc"
)

// ParseTransaction takes a raw JSON-encoded transaction and its signature,
// unmarshals it, and then extracts both token instructions and log-based events.
// It stores decoded events into the database.
func ParseTransaction(ctx context.Context, db *storage.Gorm, rawTx []byte, sig string) error {
	var tx rpc.GetTransactionResult

	// Unmarshal the raw JSON transaction
	if err := json.Unmarshal(rawTx, &tx); err != nil {
		return fmt.Errorf("unmarshal tx JSON: %w", err)
	}

	// If there are no logs, nothing to parse
	if tx.Meta == nil || tx.Meta.LogMessages == nil {
		log.Warnf("Transaction %s has no logs", sig)
		return nil
	}

	// Extract and save token-related events
	log.Infof("Parsing token instructions for transaction %s...", sig)
	if err := parseTokenInstructions(ctx, db, sig, &tx); err != nil {
		return fmt.Errorf("error parsing token instructions in %s: %w", sig, err)
	}

	// Extract and save Borsh log-based events
	log.Infof("Parsing logs for transaction %s...", sig)
	if err := parseLogs(ctx, db, sig, &tx); err != nil {
		return fmt.Errorf("error parsing logs in %s: %w", sig, err)
	}

	return nil
}

// ParseSavedTransactions iterates over all transactions that are stored in the database
// but have not yet been parsed (i.e., their `parsed` flag is false).
// It loads the raw transaction, parses it, and updates the `parsed` flag.
func ParseSavedTransactions(ctx context.Context, db *storage.Gorm, resume bool) error {
	// Get list of transactions with no parsed events
	signatures, err := db.GetOrderedNoParsedSignatures(ctx, resume)
	if err != nil {
		return err
	}
	log.Infof("Found %d signatures to parse", len(signatures))

	for _, sig := range signatures {
		// Retrieve raw transaction payload
		rawTx, err := db.GetRawTransaction(ctx, sig)
		if err != nil {
			return fmt.Errorf("failed to get raw transaction: %w", err)
		}

		// Parse and extract events from the transaction
		if err := ParseTransaction(ctx, db, rawTx, sig); err != nil {
			return fmt.Errorf("failed to parse transaction %s: %w", sig, err)
		}

		// Mark transaction as parsed in DB
		if err := db.MarkParsed(ctx, sig); err != nil {
			log.Errorf("failed to mark parsed %s: %v", sig, err)
		}
	}
	return nil
}
