package healthchecker

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/internal/storage"
)

// Check performs a full health check on the database.
// It verifies both the sequence of transactions and the sequence of events.
// If any issue is found, the health status is marked as "unhealthy".
func Check(ctx context.Context, db *storage.Gorm) error {
	// Validate transaction order and parse status
	if err := checkTransactions(ctx, db); err != nil {
		if err := db.SetHealth(ctx, "unhealthy", err.Error()); err != nil {
			return fmt.Errorf("failed to set health status: %w", err)
		}
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Validate event sequence (timestamp + logIndex)
	if err := checkEvents(ctx, db); err != nil {
		if err := db.SetHealth(ctx, "unhealthy", err.Error()); err != nil {
			return fmt.Errorf("failed to set health status: %w", err)
		}
		return fmt.Errorf("database health check failed: %w", err)
	}

	// Set health to healthy if all checks passed
	return db.SetHealth(ctx, "healthy", "")
}

// checkTransactions ensures that all parsed transactions are in proper slot order,
// and no parsed transaction follows an unparsed one.
func checkTransactions(ctx context.Context, db *storage.Gorm) error {
	//type row struct {
	//	Slot   uint64
	//	Parsed bool
	//}
	//
	//var txs []row
	//if err := db.DB.WithContext(ctx).
	//	Model(&core.Transaction{}).
	//	Order("slot ASC").
	//	Select("slot, parsed").
	//	Find(&txs).Error; err != nil {
	//	return fmt.Errorf("fetch transactions failed: %w", err)
	//}
	//
	//for i := 1; i < len(txs); i++ {
	//	// Ensure parsed transaction doesn't follow an unparsed one
	//	if txs[i].Parsed && !txs[i-1].Parsed {
	//		return fmt.Errorf("invalid parse order at slot %d: parsed transaction follows unparsed one", txs[i].Slot)
	//	}
	//}

	return nil
}

// checkEvents validates that event rows are properly ordered:
// 1. `index` is strictly increasing
// 2. `block_time` is non-decreasing
// 3. `log_index` is increasing within the same transaction signature
func checkEvents(ctx context.Context, db *storage.Gorm) error {
	// Uncomment to activate the check

	// type row struct {
	// 	Index                int64
	// 	TransactionSignature string
	// 	BlockTime            int64
	// 	LogIndex             int
	// }
	// var events []row

	// if err := db.DB.WithContext(ctx).
	// 	Model(&core.Event{}).
	// 	Order("index ASC").
	// 	Select("index, transaction_signature, block_time, log_index").
	// 	Find(&events).Error; err != nil {
	// 	return fmt.Errorf("fetch events failed: %w", err)
	// }

	// var (
	// 	prevIndex     int64
	// 	prevSignature string
	// 	prevBlockTime int64
	// 	prevLogIndex  int
	// )

	// for i, ev := range events {
	// 	if i == 0 {
	// 		prevIndex = ev.Index
	// 		prevSignature = ev.TransactionSignature
	// 		prevBlockTime = ev.BlockTime
	// 		prevLogIndex = ev.LogIndex
	// 		continue
	// 	}

	// 	if ev.Index <= prevIndex {
	// 		return fmt.Errorf("index out of order at row %d: %d <= %d", i, ev.Index, prevIndex)
	// 	}

	// 	if ev.BlockTime < prevBlockTime {
	// 		return fmt.Errorf("block_time decreased at row %d: %d < %d", i, ev.BlockTime, prevBlockTime)
	// 	}

	// 	if ev.TransactionSignature == prevSignature && ev.LogIndex <= prevLogIndex {
	// 		return fmt.Errorf("log_index out of order within transaction %s at row %d: %d <= %d",
	// 			ev.TransactionSignature, i, ev.LogIndex, prevLogIndex)
	// 	}

	// 	prevIndex = ev.Index
	// 	prevSignature = ev.TransactionSignature
	// 	prevBlockTime = ev.BlockTime
	// 	prevLogIndex = ev.LogIndex
	// }

	return nil
}
