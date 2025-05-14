package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/gagliardetto/solana-go/rpc"
)

func parseTransactions(ctx context.Context, db *storage.Gorm, sig string) error {
	rawTx, err := db.GetRawTransaction(ctx, sig)
	if err != nil {
		return fmt.Errorf("failed to get raw transaction: %w", err)
	}
	if err := ParseTransaction(ctx, db, rawTx, sig); err != nil {
		log.Errorf("failed to process tx %s: %v", sig, err)
		return err
	}
	if err := db.MarkParsed(ctx, sig); err != nil {
		log.Errorf("failed to mark parsed %s: %v", sig, err)
	}
	return nil
}

func ParseTransaction(ctx context.Context, db *storage.Gorm, rawTx []byte, sig string) error {
	// TODO: винести в окрему функцію, дял того аби парсити логи чи інструкції просто передавати сигнатуру зберігати також в ній
	var tx rpc.GetTransactionResult
	if err := json.Unmarshal(rawTx, &tx); err != nil {
		return fmt.Errorf("unmarshal tx JSON: %w", err)
	}
	if tx.Meta == nil || tx.Meta.LogMessages == nil {
		return nil
	}

	log.Infof("Parsing token instructions for transaction %s...", sig)
	if err := parseTokenInstructions(ctx, db, sig, &tx); err != nil {
		log.Errorf("Error parsing token instructions in %s: %v", sig, err)
	}

	log.Infof("Parsing logs for transaction %s...", sig)
	if err := parseLogs(ctx, db, sig, &tx); err != nil {
		log.Errorf("Error parsing logs in %s: %v", sig, err)
	}

	return nil
}

func Parse(ctx context.Context, db *storage.Gorm) error {
	signatures, err := db.GetOrderedNoParsedSignatures(ctx)
	if err != nil {
		return err
	}
	for _, sig := range signatures {
		//log.Infof("Parsing transaction %s", sig)
		if err := parseTransactions(ctx, db, sig); err != nil {
			// TODO mark db as unhealthy
			return err
		}
	}
	return nil
}

func Start(ctx context.Context, db *storage.Gorm) {
	// TODO lock in db
	if err := Parse(ctx, db); err != nil {
		log.Fatalf("Failed to parse logs: %v", err)
	}
}
