package strategy

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/subgraph/events"

	"gorm.io/gorm"
)

func Init(ctx context.Context, db *gorm.DB, ev events.StrategyInitEvent) error {
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[Init] failed to load strategy: %w", err)
	}

	strategy.StrategyType = ev.StrategyType
	strategy.DepositLimit = ev.DepositLimit.String()
	strategy.DepositPeriodEnds = ev.DepositPeriodEnds.String()
	strategy.LockPeriodEnds = ev.LockPeriodEnds.String()
	strategy.VaultID = ev.Vault.String()
	strategy.Removed = false

	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[Init] failed to save strategy: %w", err)
	}
	return nil
}
