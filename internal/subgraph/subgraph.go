package subgraph

import (
	"context"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/Tsisar/solana-indexer/internal/storage/model/core"
	"github.com/Tsisar/solana-indexer/internal/subgraph/aggregator"
	"github.com/Tsisar/solana-indexer/internal/subgraph/maping"
)

func MapEvent(ctx context.Context, db *storage.Gorm, event core.Event) {
	if err := maping.Event(ctx, db.DB, event); err != nil {
		if err := maping.Error(ctx, db.DB, err); err != nil {
			log.Errorf("Failed to map error: %v", err)
		}
		log.Fatalf("Failed to map event: %v", err)
	}
}

func MapInstruction(ctx context.Context, db *storage.Gorm, event core.Event) {
	if err := maping.Instruction(ctx, db.DB, event); err != nil {
		if err := maping.Error(ctx, db.DB, err); err != nil {
			log.Errorf("Failed to map error: %v", err)
		}
		log.Fatalf("Failed to map instruction: %v", err)
	}
}

func MapMetadata(ctx context.Context, db *storage.Gorm, signature string, slot uint64, blockTime int64) {
	if err := maping.Metadata(ctx, db.DB, signature, slot, blockTime); err != nil {
		if err := maping.Error(ctx, db.DB, err); err != nil {
			log.Errorf("Failed to map error: %v", err)
		}
		log.Fatalf("Failed to map metadata: %v", err)
	}

}

func RunAggregator(ctx context.Context, db *storage.Gorm) {
	if err := aggregator.Start(ctx, db.DB); err != nil {
		if err := maping.Error(ctx, db.DB, err); err != nil {
			log.Errorf("Failed to map error: %v", err)
		}
		log.Errorf("Failed to run aggregator: %v", err)
	}
}

func MapError(ctx context.Context, db *storage.Gorm, err error) {
	if err := maping.Error(ctx, db.DB, err); err != nil {
		log.Fatalf("Failed to map error: %v", err)
	}
}
