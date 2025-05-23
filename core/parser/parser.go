package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/utils"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/subgraph"
	"github.com/gagliardetto/solana-go/rpc"
)

// Start processes all unparsed transactions from the DB
// and then continues parsing from the real-time stream.
// Returns error if any transaction fails to parse.
func Start(ctx context.Context, db *storage.Gorm, resume bool, done chan struct{}, in <-chan string) error {
	// Load list of signatures that are not parsed
	signatures, err := db.GetOrderedNoParsedSignatures(ctx, resume)
	if err != nil {
		return fmt.Errorf("[parser] failed to load signatures to parse: %w", err)
	}
	log.Infof("[parser] Found %d transactions to parse", len(signatures))

	for _, sig := range signatures {
		if err := parseOneTransaction(ctx, db, resume, sig); err != nil {
			return fmt.Errorf("[parser] failed to parse transaction %s: %w", sig, err)
		}
	}
	done <- struct{}{}

	log.Infof("[parser] done with DB, switching to real-time stream...")

	// Switch to processing incoming signatures from WebSocket stream
	for {
		select {
		case <-ctx.Done():
			log.Debugf("[parser] context cancelled")
			return nil
		case sig := <-in:
			if err := parseOneTransaction(ctx, db, true, sig); err != nil {
				return fmt.Errorf("[parser] failed to parse real-time transaction %s: %w", sig, err)
			}
		}
	}
}

// parseOneTransaction coordinates parsing of one transaction from DB by signature.
func parseOneTransaction(ctx context.Context, db *storage.Gorm, resume bool, sig string) error {
	if resume {
		parsed, err := db.IsParsed(ctx, sig)
		if err != nil {
			return fmt.Errorf("[parser] failed to check if %s is parsed: %w", sig, err)
		}
		if parsed {
			log.Warnf("[parser] Transaction %s already parsed, skipping...", sig)
			return nil
		}
	}

	// Retrieve raw transaction from DB
	rawTx, err := db.GetRawTransaction(ctx, sig)
	if err != nil {
		return fmt.Errorf("[parser] failed to get raw transaction %s: %w", sig, err)
	}

	// Parse and store token instructions and logs
	if err := parseTransaction(ctx, db, rawTx, sig); err != nil {
		return fmt.Errorf("[parser] failed to parse transaction %s: %w", sig, err)
	}

	// Mark transaction as parsed in DB
	if err := db.MarkParsed(ctx, sig); err != nil {
		return fmt.Errorf("[parser] failed to mark %s as parsed: %w", sig, err)
	}

	return nil
}

// parseTransaction unmarshals the JSON payload and extracts events and instructions.
func parseTransaction(ctx context.Context, db *storage.Gorm, rawTx []byte, sig string) error {
	var tx rpc.GetTransactionResult

	// Decode JSON
	if err := json.Unmarshal(rawTx, &tx); err != nil {
		return fmt.Errorf("[parser] unmarshal tx JSON: %w", err)
	}

	if tx.Meta == nil || tx.Meta.LogMessages == nil {
		log.Warnf("[parser] Transaction %s has no logs", sig)
		return nil
	}

	log.Infof("[parser] Parsing instructions for %s", sig)
	if err := parseTokenInstructions(ctx, db, sig, &tx); err != nil {
		return fmt.Errorf("[parser] error parsing instructions in %s: %w", sig, err)
	}

	log.Infof("[parser] Parsing logs for %s", sig)
	if err := parseLogs(ctx, db, sig, &tx); err != nil {
		return fmt.Errorf("[parser] error parsing logs in %s: %w", sig, err)
	}

	subgraph.MapMetadata(ctx, db, sig, tx.Slot, utils.BlockTime(tx.BlockTime))

	return nil
}
