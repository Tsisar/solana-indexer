package maping

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/Tsisar/solana-indexer/subgraph/events"
	"github.com/Tsisar/solana-indexer/subgraph/library/strategy"
	"github.com/Tsisar/solana-indexer/subgraph/library/vault"
	"gorm.io/gorm"
)

func mapStrategyReportedEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] StrategyReportedEvent: %s", event.TransactionSignature)
	var ev events.StrategyReportedEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode StrategyReportedEvent: %w", err)
	}

	transaction := events.NewTransaction(event)

	if err := vault.StrategyReported(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[maping] failed to report strategy: %w", err)

	}

	return nil
}

func mapUpdatedCurrentDebtForStrategyEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] UpdatedCurrentDebtForStrategyEvent: %s", event.TransactionSignature)
	var ev events.UpdatedCurrentDebtForStrategyEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode UpdatedCurrentDebtForStrategyEvent: %w", err)
	}
	if err := strategy.UpdateCurrentDebt(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to update current debt: %w", err)
	}
	return nil
}

func mapVaultAddStrategyEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultAddStrategyEvent: %s", event.TransactionSignature)
	var ev events.VaultAddStrategyEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultAddStrategyEvent: %w", err)
	}

	transaction := events.NewTransaction(event)

	if err := vault.AddStrategy(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[maping] failed to add strategy: %w", err)
	}
	return nil
}

func mapVaultDepositEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultDepositEvent: %s", event.TransactionSignature)
	var ev events.VaultDepositEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultDepositEvent: %w", err)
	}

	transaction := events.NewTransaction(event)

	if err := vault.Deposit(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[maping] failed to deposit: %w", err)
	}
	return nil
}

func mapVaultEmergencyWithdrawEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultEmergencyWithdrawEvent: %s", event.TransactionSignature)
	var ev events.VaultEmergencyWithdrawEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultEmergencyWithdrawEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultInitEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultInitEvent: %s", event.TransactionSignature)
	var ev events.VaultInitEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultInitEvent: %w", err)
	}

	transaction := events.NewTransaction(event)

	if err := vault.Init(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[maping] failed to init vault: %w", err)
	}

	return nil
}

func mapVaultRemoveStrategyEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultRemoveStrategyEvent: %s", event.TransactionSignature)
	var ev events.VaultRemoveStrategyEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode VaultRemoveStrategyEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultShutDownEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultShutDownEvent: %s", event.TransactionSignature)
	var ev events.VaultShutDownEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultShutDownEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultUpdateAccountantEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateAccountantEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateAccountantEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateAccountantEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultUpdateDepositLimitEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateDepositLimitEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateDepositLimitEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateDepositLimitEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultUpdateDirectWithdrawEnabledEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateDirectWithdrawEnabledEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateDirectWithdrawEnabledEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateDirectWithdrawEnabledEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultUpdateMinTotalIdleEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateMinTotalIdleEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateMinTotalIdleEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateMinTotalIdleEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultUpdateMinUserDepositEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateMinUserDepositEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateMinUserDepositEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateMinUserDepositEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultUpdateProfitMaxUnlockTimeEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateProfitMaxUnlockTimeEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateProfitMaxUnlockTimeEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateProfitMaxUnlockTimeEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultUpdateUserDepositLimitEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateUserDepositLimitEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateUserDepositLimitEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateUserDepositLimitEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultUpdateWhitelistedOnlyEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateWhitelistedOnlyEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateWhitelistedOnlyEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateWhitelistedOnlyEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapVaultWithdrawlEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultWithdrawlEvent: %s", event.TransactionSignature)
	var ev events.VaultWithdrawlEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultWithdrawlEvent: %w", err)
	}

	transaction := events.NewTransaction(event)

	if err := vault.Withdraw(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[maping] failed to withdraw: %w", err)
	}
	return nil
}

func mapWhitelistUpdatedEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] WhitelistUpdatedEvent: %s", event.TransactionSignature)
	var ev events.WhitelistUpdatedEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode WhitelistUpdatedEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapWithdrawalRequestCanceledEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] WithdrawalRequestCanceledEvent: %s", event.TransactionSignature)
	var ev events.WithdrawalRequestCanceledEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode WithdrawalRequestCanceledEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapWithdrawalRequestFulfilledEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] WithdrawalRequestFulfilledEvent: %s", event.TransactionSignature)
	var ev events.WithdrawalRequestFulfilledEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode WithdrawalRequestFulfilledEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapWithdrawalRequestedEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] WithdrawalRequestedEvent: %s", event.TransactionSignature)
	var ev events.WithdrawalRequestedEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode WithdrawalRequestedEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}
