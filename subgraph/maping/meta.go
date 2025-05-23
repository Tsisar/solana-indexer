package maping

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/core/config"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"gorm.io/gorm"
)

func updateMeta(ctx context.Context, db *gorm.DB, signature string, slot uint64, blockTime int64) error {
	meta := subgraph.Meta{
		ID:                1,
		Deployment:        fmt.Sprintf("solana-indexer %s", config.App.Version),
		HasIndexingErrors: false,
		BlockID:           1,
		Block: &subgraph.BlockInfo{
			ID:         1,
			Hash:       signature,
			Number:     slot,
			ParentHash: signature,
			Timestamp:  blockTime,
		},
	}

	if err := meta.Block.Save(ctx, db); err != nil {
		return fmt.Errorf("failed to save block: %w", err)
	}

	if err := meta.Save(ctx, db); err != nil {
		return fmt.Errorf("failed to save meta: %w", err)
	}
	return nil
}

func Error(ctx context.Context, db *gorm.DB, err error) error {
	meta := subgraph.Meta{
		ID:                1,
		Deployment:        "solana-indexer",
		HasIndexingErrors: true,
		ErrorMessage:      err.Error(),
		Block: &subgraph.BlockInfo{
			ID: 1,
		},
	}

	if err := meta.Block.Save(ctx, db); err != nil {
		return fmt.Errorf("failed to save block: %w", err)
	}

	if err := meta.Save(ctx, db); err != nil {
		return fmt.Errorf("failed to save meta: %w", err)
	}
	return nil
}
