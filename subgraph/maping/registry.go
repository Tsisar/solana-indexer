package maping

import (
	"context"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"gorm.io/gorm"
)

type EventMapper func(ctx context.Context, db *gorm.DB, event core.Event) error

var registry = map[string]EventMapper{
	"EntryFeeUpdatedEvent":                  mapEntryFeeUpdatedEvent,       //Done
	"PerformanceFeeUpdatedEvent":            mapPerformanceFeeUpdatedEvent, //Done
	"RedemptionFeeUpdatedEvent":             mapRedemptionFeeUpdatedEvent,  //Done
	"DepositLimitSetEvent":                  mapDepositLimitSetEvent,
	"EmergencyWithdrawEvent":                mapEmergencyWithdrawEvent,
	"FundManagerDeployFundsEvent":           mapFundManagerDeployFundsEvent,
	"FundManagerEmergencyWithdrawEvent":     mapFundManagerEmergencyWithdrawEvent,
	"FundManagerFreeFundsEvent":             mapFundManagerFreeFundsEvent,
	"FundManagerHarvestAndReportEvent":      mapFundManagerHarvestAndReportEvent,
	"FundManagerStrategyStateUpdateEvent":   mapFundManagerStrategyStateUpdateEvent,
	"HarvestAndReportDTFEvent":              mapHarvestAndReportDTFEvent,
	"MinDeployAmountSetEvent":               mapMinDeployAmountSetEvent,
	"OrcaAfterSwapEvent":                    mapOrcaAfterSwapEvent,
	"OrcaInitEvent":                         mapOrcaInitEvent,
	"SetPerformanceFeeEvent":                mapSetPerformanceFeeEvent,
	"StrategyDeployFundsEvent":              mapStrategyDeployFundsEvent,
	"StrategyDepositEvent":                  mapStrategyDepositEvent,
	"StrategyFreeFundsEvent":                mapStrategyFreeFundsEvent,
	"StrategyInitEvent":                     mapStrategyInitEvent, //Done
	"StrategyReallocEvent":                  mapStrategyReallocEvent,
	"StrategyShutdownEvent":                 mapStrategyShutdownEvent,
	"StrategyWithdrawEvent":                 mapStrategyWithdrawEvent,
	"TotalInvestedUpdatedEvent":             mapTotalInvestedUpdatedEvent,
	"StrategyReportedEvent":                 mapStrategyReportedEvent,
	"UpdatedCurrentDebtForStrategyEvent":    mapUpdatedCurrentDebtForStrategyEvent,
	"VaultAddStrategyEvent":                 mapVaultAddStrategyEvent, //Done
	"VaultDepositEvent":                     mapVaultDepositEvent,
	"VaultEmergencyWithdrawEvent":           mapVaultEmergencyWithdrawEvent,
	"VaultInitEvent":                        mapVaultInitEvent, //Done
	"VaultRemoveStrategyEvent":              mapVaultRemoveStrategyEvent,
	"VaultShutDownEvent":                    mapVaultShutDownEvent,
	"VaultUpdateAccountantEvent":            mapVaultUpdateAccountantEvent,
	"VaultUpdateDepositLimitEvent":          mapVaultUpdateDepositLimitEvent,
	"VaultUpdateDirectWithdrawEnabledEvent": mapVaultUpdateDirectWithdrawEnabledEvent,
	"VaultUpdateMinTotalIdleEvent":          mapVaultUpdateMinTotalIdleEvent,
	"VaultUpdateMinUserDepositEvent":        mapVaultUpdateMinUserDepositEvent,
	"VaultUpdateProfitMaxUnlockTimeEvent":   mapVaultUpdateProfitMaxUnlockTimeEvent,
	"VaultUpdateUserDepositLimitEvent":      mapVaultUpdateUserDepositLimitEvent,
	"VaultUpdateWhitelistedOnlyEvent":       mapVaultUpdateWhitelistedOnlyEvent,
	"VaultWithdrawlEvent":                   mapVaultWithdrawlEvent,
	"WhitelistUpdatedEvent":                 mapWhitelistUpdatedEvent,
	"WithdrawalRequestCanceledEvent":        mapWithdrawalRequestCanceledEvent,
	"WithdrawalRequestFulfilledEvent":       mapWithdrawalRequestFulfilledEvent,
	"WithdrawalRequestedEvent":              mapWithdrawalRequestedEvent,
}

func mapEvents(ctx context.Context, db *gorm.DB, event core.Event) error {
	if handler, ok := registry[event.Name]; ok {
		return handler(ctx, db, event)
	}
	log.Warnf("No mapping implemented for event: %s", event.Name)
	return nil
}
