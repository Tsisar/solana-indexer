package maping

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/storage/model/core"
	"github.com/Tsisar/solana-indexer/internal/subgraph/events"
	"github.com/Tsisar/solana-indexer/internal/subgraph/library/shareToken"
	"github.com/Tsisar/solana-indexer/internal/subgraph/library/strategy"
	"github.com/Tsisar/solana-indexer/internal/subgraph/library/vault"
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
		return fmt.Errorf("[maping] failed to decode VaultRemoveStrategyEvent: %w", err)
	}
	if err := strategy.Remove(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to remove strategy: %w", err)
	}
	return nil
}

func mapVaultShutDownEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultShutDownEvent: %s", event.TransactionSignature)
	var ev events.VaultShutDownEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultShutDownEvent: %w", err)
	}

	if err := vault.ShutDown(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to shut down vault: %w", err)
	}
	return nil
}

func mapVaultUpdateAccountantEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateAccountantEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateAccountantEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateAccountantEvent: %w", err)
	}
	if err := vault.UpdateAccountant(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to update accountant: %w", err)
	}
	return nil
}

func mapVaultUpdateDepositLimitEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateDepositLimitEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateDepositLimitEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateDepositLimitEvent: %w", err)
	}

	if err := vault.UpdateDepositLimit(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to update deposit limit: %w", err)
	}
	return nil
}

func mapVaultUpdateDirectWithdrawEnabledEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateDirectWithdrawEnabledEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateDirectWithdrawEnabledEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateDirectWithdrawEnabledEvent: %w", err)
	}
	if err := vault.UpdateDirectWithdrawEnabled(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to update direct withdraw enabled: %w", err)
	}
	return nil
}

func mapVaultUpdateMinTotalIdleEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateMinTotalIdleEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateMinTotalIdleEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateMinTotalIdleEvent: %w", err)
	}
	if err := vault.UpdateMinTotalIdle(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to update min total idle: %w", err)
	}
	return nil
}

func mapVaultUpdateMinUserDepositEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateMinUserDepositEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateMinUserDepositEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateMinUserDepositEvent: %w", err)
	}
	if err := vault.UpdateMinUserDeposit(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to update min user deposit: %w", err)
	}
	return nil
}

func mapVaultUpdateProfitMaxUnlockTimeEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateProfitMaxUnlockTimeEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateProfitMaxUnlockTimeEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateProfitMaxUnlockTimeEvent: %w", err)
	}
	if err := vault.UpdateProfitMaxUnlockTime(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to update profit max unlock time: %w", err)
	}
	return nil
}

func mapVaultUpdateUserDepositLimitEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateUserDepositLimitEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateUserDepositLimitEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateUserDepositLimitEvent: %w", err)
	}
	if err := vault.UpdateUserDepositLimit(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to update user deposit limit: %w", err)
	}
	return nil
}

func mapVaultUpdateWhitelistedOnlyEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] VaultUpdateWhitelistedOnlyEvent: %s", event.TransactionSignature)
	var ev events.VaultUpdateWhitelistedOnlyEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode VaultUpdateWhitelistedOnlyEvent: %w", err)
	}
	if err := vault.UpdateWhiteListOnly(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to update whitelist only: %w", err)
	}
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
	if err := vault.WithdrawalRequestCanceled(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to cancel withdrawal request: %w", err)
	}
	return nil
}

func mapWithdrawalRequestFulfilledEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] WithdrawalRequestFulfilledEvent: %s", event.TransactionSignature)
	var ev events.WithdrawalRequestFulfilledEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode WithdrawalRequestFulfilledEvent: %w", err)
	}
	if err := vault.WithdrawalRequestFulfilled(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to fulfill withdrawal request: %w", err)
	}
	return nil
}

func mapWithdrawalRequestedEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] WithdrawalRequestedEvent: %s", event.TransactionSignature)
	var ev events.WithdrawalRequestedEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode WithdrawalRequestedEvent: %w", err)
	}
	if err := vault.WithdrawalRequested(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to process withdrawal request: %w", err)
	}
	return nil
}

func mapMintToInstruction(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] MintToInstruction: %s", event.TransactionSignature)
	var ev events.MintToInstruction
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode MintToInstruction: %w", err)
	}
	transaction := events.NewTransaction(event)

	if err := shareToken.Mint(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[maping] failed to mint: %w", err)
	}
	return nil
}

func mapBurnInstruction(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] BurnInstruction: %s", event.TransactionSignature)
	var ev events.BurnInstruction
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode BurnInstruction: %w", err)
	}
	transaction := events.NewTransaction(event)

	if err := shareToken.Burn(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[maping] failed to burn: %w", err)
	}
	return nil
}

func mapTransferInstruction(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] TransferInstruction: %s", event.TransactionSignature)
	var ev events.TransferInstruction
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode TransferInstruction: %w", err)
	}
	transaction := events.NewTransaction(event)

	if err := shareToken.Transfer(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[maping] failed to transfer: %w", err)
	}
	return nil
}

func mapInitializeAccount3Instruction(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[maping] InitializeAccountInstruction: %s", event.TransactionSignature)
	var ev events.InitializeAccountInstruction
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[maping] failed to decode InitializeAccountInstruction: %w", err)
	}

	if err := shareToken.InitializeAccount(ctx, db, ev); err != nil {
		return fmt.Errorf("[maping] failed to initialize account: %w", err)
	}
	return nil
}
