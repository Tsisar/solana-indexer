package maping

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/Tsisar/solana-indexer/subgraph/events"
	"github.com/Tsisar/solana-indexer/subgraph/library/strategy"
	"gorm.io/gorm"
)

func mapDepositLimitSetEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] DepositLimitSetEvent: %s", event.TransactionSignature)
	var ev events.DepositLimitSetEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("[mapping] failed to decode DepositLimitSetEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapEmergencyWithdrawEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] EmergencyWithdrawEvent: %s", event.TransactionSignature)
	var ev events.EmergencyWithdrawEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode EmergencyWithdrawEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapFundManagerDeployFundsEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] FundManagerDeployFundsEvent: %s", event.TransactionSignature)
	var ev events.FundManagerDeployFundsEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode FundManagerDeployFundsEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapFundManagerEmergencyWithdrawEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] FundManagerEmergencyWithdrawEvent: %s", event.TransactionSignature)
	var ev events.FundManagerEmergencyWithdrawEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode FundManagerEmergencyWithdrawEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapFundManagerFreeFundsEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] FundManagerFreeFundsEvent: %s", event.TransactionSignature)
	var ev events.FundManagerFreeFundsEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode FundManagerFreeFundsEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapFundManagerHarvestAndReportEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] FundManagerHarvestAndReportEvent: %s", event.TransactionSignature)
	var ev events.FundManagerHarvestAndReportEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode FundManagerHarvestAndReportEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapFundManagerStrategyStateUpdateEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] FundManagerStrategyStateUpdateEvent: %s", event.TransactionSignature)
	var ev events.FundManagerStrategyStateUpdateEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode FundManagerStrategyStateUpdateEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapHarvestAndReportDTFEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] HarvestAndReportDTFEvent: %s", event.TransactionSignature)
	var ev events.HarvestAndReportDTFEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode HarvestAndReportDTFEvent: %w", err)
	}
	if err := strategy.UpdateDTFReport(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to update DTF report: %w", err)
	}
	return nil
}

func mapMinDeployAmountSetEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] MinDeployAmountSetEvent: %s", event.TransactionSignature)
	var ev events.MinDeployAmountSetEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode MinDeployAmountSetEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapOrcaAfterSwapEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] OrcaAfterSwapEvent: %s", event.TransactionSignature)
	var ev events.OrcaAfterSwapEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode OrcaAfterSwapEvent: %w", err)
	}
	if err := strategy.AfterOrcaSwap(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to process OrcaAfterSwapEvent: %w", err)
	}
	return nil
}

func mapOrcaInitEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] OrcaInitEvent: %s", event.TransactionSignature)
	var ev events.OrcaInitEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode OrcaInitEvent: %w", err)
	}
	if err := strategy.InitOrca(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to process OrcaInitEvent: %w", err)
	}
	return nil
}

func mapSetPerformanceFeeEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] SetPerformanceFeeEvent: %s", event.TransactionSignature)
	var ev events.SetPerformanceFeeEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode SetPerformanceFeeEvent: %w", err)
	}
	if err := strategy.UpdatePerformanceFee(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to update performance fee: %w", err)
	}
	return nil
}

func mapStrategyDeployFundsEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] StrategyDeployFundsEvent: %s", event.TransactionSignature)
	var ev events.StrategyDeployFundsEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode StrategyDeployFundsEvent: %w", err)
	}

	if err := strategy.DeployFunds(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to deploy funds: %w", err)
	}
	return nil
}

func mapStrategyDepositEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] StrategyDepositEvent: %s", event.TransactionSignature)
	var ev events.StrategyDepositEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode StrategyDepositEvent: %w", err)
	}

	if err := strategy.Deposit(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to deposit strategy: %w", err)
	}
	return nil
}

func mapStrategyFreeFundsEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] StrategyFreeFundsEvent: %s", event.TransactionSignature)
	var ev events.StrategyFreeFundsEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode StrategyFreeFundsEvent: %w", err)
	}
	if err := strategy.FreeFunds(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to free funds: %w", err)
	}
	return nil
}

func mapStrategyInitEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] StrategyInitEvent: %s", event.TransactionSignature)
	var ev events.StrategyInitEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode StrategyInitEvent: %w", err)
	}
	if err := strategy.Init(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to initialize strategy: %w", err)
	}
	return nil
}

func mapStrategyReallocEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] StrategyReallocEvent: %s", event.TransactionSignature)
	var ev events.StrategyReallocEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode StrategyReallocEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapStrategyShutdownEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] StrategyShutdownEvent: %s", event.TransactionSignature)
	var ev events.StrategyShutdownEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode StrategyShutdownEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}

func mapStrategyWithdrawEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] StrategyWithdrawEvent: %s", event.TransactionSignature)
	var ev events.StrategyWithdrawEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode StrategyWithdrawEvent: %w", err)
	}

	if err := strategy.Withdraw(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to withdraw strategy: %w", err)
	}
	return nil
}

func mapTotalInvestedUpdatedEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] TotalInvestedUpdatedEvent: %s", event.TransactionSignature)
	var ev events.TotalInvestedUpdatedEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode TotalInvestedUpdatedEvent: %w", err)
	}
	// TODO: implement mapping logic
	return nil
}
