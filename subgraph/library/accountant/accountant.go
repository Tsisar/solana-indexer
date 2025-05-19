package accountant

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/subgraph/events"
	"gorm.io/gorm"
)

func Init(ctx context.Context, db *gorm.DB, id string) (*subgraph.Accountant, error) {
	accountant := subgraph.Accountant{ID: id}
	ok, err := accountant.Load(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("[Init] failed to load accountant: %w", err)
	}
	if !ok {
		if err := accountant.Save(ctx, db); err != nil {
			return nil, fmt.Errorf("[Init] failed to save accountant: %w", err)
		}
	}
	return &accountant, nil
}

func SetRedemptionFee(ctx context.Context, db *gorm.DB, ev events.RedemptionFeeUpdatedEvent) error {
	accountant := subgraph.Accountant{ID: ev.AccountantKey.String()}
	if _, err := accountant.Load(ctx, db); err != nil {
		return fmt.Errorf("[SetRedemptionFee] failed to load accountant: %w", err)
	}
	accountant.RedemptionFee = ev.RedemptionFee
	if err := accountant.Save(ctx, db); err != nil {
		return fmt.Errorf("[SetRedemptionFee] failed to save accountant: %w", err)
	}
	return nil
}

func SetPerformanceFee(ctx context.Context, db *gorm.DB, ev events.PerformanceFeeUpdatedEvent) error {
	accountant := subgraph.Accountant{ID: ev.AccountantKey.String()}
	if _, err := accountant.Load(ctx, db); err != nil {
		return fmt.Errorf("[SetPerformanceFee] failed to load accountant: %w", err)
	}
	accountant.PerformanceFees = ev.PerformanceFee
	if err := accountant.Save(ctx, db); err != nil {
		return fmt.Errorf("[SetPerformanceFee] failed to save accountant: %w", err)
	}
	return nil
}

func SetEntryFee(ctx context.Context, db *gorm.DB, ev events.EntryFeeUpdatedEvent) error {
	accountant := subgraph.Accountant{ID: ev.AccountantKey.String()}
	if _, err := accountant.Load(ctx, db); err != nil {
		return fmt.Errorf("[SetEntryFee] failed to load accountant: %w", err)
	}
	accountant.EntryFee = ev.EntryFee
	if err := accountant.Save(ctx, db); err != nil {
		return fmt.Errorf("[SetEntryFee] failed to save accountant: %w", err)
	}
	return nil
}
