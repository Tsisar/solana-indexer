package vault

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/subgraph/events"
	"github.com/Tsisar/solana-indexer/subgraph/library/accountant"
	"github.com/Tsisar/solana-indexer/subgraph/library/token"

	"gorm.io/gorm"
)

func Init(ctx context.Context, db *gorm.DB, ev events.VaultInitEvent, transaction events.Transaction) error {
	var err error
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err = vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[Init] failed to load vault: %w", err)
	}

	vault.Token, err = token.UpsertUnderlyingToken(ctx, db, ev)
	if err != nil {
		return fmt.Errorf("[Init] failed to get or create underlying token entity: %v", err)
	}

	vault.ShareToken, err = token.UpsertShareToken(ctx, db, ev)
	if err != nil {
		return fmt.Errorf("[Init] failed to get or create share token entity: %v", err)
	}

	acc, err := accountant.Init(ctx, db, ev.Accountant.String())
	if err != nil {
		return fmt.Errorf("[Init] failed to upsert accountant: %w", err)
	}

	vault.Shutdown = false
	vault.Activation = transaction.Timestamp
	vault.AccountantID = acc.ID
	vault.MinUserDeposit = ev.MinUserDeposit
	vault.KycVerifiedOnly = ev.KycVerifiedOnly
	vault.DirectDepositEnabled = ev.DirectDepositEnabled
	vault.WhitelistedOnly = ev.WhitelistedOnly
	vault.ProfitMaxUnlockTime = ev.ProfitMaxUnlockTime
	vault.LastUpdate = transaction.Timestamp
	vault.MinTotalIdle = ev.MinimumTotalIdle
	vault.DirectWithdrawEnabled = ev.DirectWithdrawEnabled
	vault.UserDepositLimit = ev.UserDepositLimit

	if err = vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[Init] failed to save vault: %w", err)
	}

	return nil
}

func AddStrategy(ctx context.Context, db *gorm.DB, ev events.VaultAddStrategyEvent, transaction events.Transaction) error {
	strategy := subgraph.Strategy{
		ID: ev.StrategyKey.String(),
	}
	ok, err := strategy.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[AddStrategy] failed to load strategy: %w", err)
	}

	strategy.MaxDebt = ev.MaxDebt
	strategy.CurrentDebt = ev.CurrentDebt
	strategy.VaultID = ev.VaultKey.String()
	strategy.Activation = transaction.Timestamp

	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[AddStrategy] failed to save strategy: %w", err)
	}
	return nil
}
