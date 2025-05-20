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
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}

	strategy.StrategyType = ev.StrategyType
	strategy.DepositLimit = ev.DepositLimit
	strategy.DepositPeriodEnds = ev.DepositPeriodEnds
	strategy.LockPeriodEnds = ev.LockPeriodEnds
	strategy.VaultID = ev.Vault.String()
	strategy.UnderlyingMint = ev.UnderlyingMint.String()
	strategy.UnderlyingTokenAcc = ev.UnderlyingTokenAcc.String()
	strategy.UnderlyingDecimals = ev.UnderlyingDecimals

	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}
	return nil
}

func Deposit(ctx context.Context, db *gorm.DB, ev events.StrategyDepositEvent) error {
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	strategy.TotalAssets = ev.TotalAssets
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}

	return nil
}

func Withdraw(ctx context.Context, db *gorm.DB, ev events.StrategyWithdrawEvent) error {
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	strategy.TotalAssets = ev.TotalAssets
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}

	return nil
}

func UpdateCurrentDebt(ctx context.Context, db *gorm.DB, ev events.UpdatedCurrentDebtForStrategyEvent) error {
	strategy := subgraph.Strategy{ID: ev.StrategyKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	strategy.CurrentDebt = ev.NewDebt
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}

	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err := vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load vault: %w", err)
	}
	vault.TotalDebt = ev.TotalDebt
	vault.TotalIdle = ev.TotalIdle
	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save vault: %w", err)
	}

	return nil
}
