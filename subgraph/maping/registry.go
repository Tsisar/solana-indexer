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
	"HarvestAndReportDTFEvent":              mapHarvestAndReportDTFEvent, //Done
	"MinDeployAmountSetEvent":               mapMinDeployAmountSetEvent,
	"OrcaAfterSwapEvent":                    mapOrcaAfterSwapEvent,       //Done
	"OrcaInitEvent":                         mapOrcaInitEvent,            //Done
	"SetPerformanceFeeEvent":                mapSetPerformanceFeeEvent,   //Done
	"StrategyDeployFundsEvent":              mapStrategyDeployFundsEvent, //Done
	"StrategyDepositEvent":                  mapStrategyDepositEvent,     //Done
	"StrategyFreeFundsEvent":                mapStrategyFreeFundsEvent,   //Done
	"StrategyInitEvent":                     mapStrategyInitEvent,        //Done
	"StrategyReallocEvent":                  mapStrategyReallocEvent,
	"StrategyShutdownEvent":                 mapStrategyShutdownEvent,
	"StrategyWithdrawEvent":                 mapStrategyWithdrawEvent, //Done
	"TotalInvestedUpdatedEvent":             mapTotalInvestedUpdatedEvent,
	"StrategyReportedEvent":                 mapStrategyReportedEvent,              //Done
	"UpdatedCurrentDebtForStrategyEvent":    mapUpdatedCurrentDebtForStrategyEvent, //Done
	"VaultAddStrategyEvent":                 mapVaultAddStrategyEvent,              //Done
	"VaultDepositEvent":                     mapVaultDepositEvent,                  //Done
	"VaultEmergencyWithdrawEvent":           mapVaultEmergencyWithdrawEvent,
	"VaultInitEvent":                        mapVaultInitEvent,                        //Done
	"VaultRemoveStrategyEvent":              mapVaultRemoveStrategyEvent,              //Done
	"VaultShutDownEvent":                    mapVaultShutDownEvent,                    //Done
	"VaultUpdateAccountantEvent":            mapVaultUpdateAccountantEvent,            //Done
	"VaultUpdateDepositLimitEvent":          mapVaultUpdateDepositLimitEvent,          //Done
	"VaultUpdateDirectWithdrawEnabledEvent": mapVaultUpdateDirectWithdrawEnabledEvent, //Done
	"VaultUpdateMinTotalIdleEvent":          mapVaultUpdateMinTotalIdleEvent,          //Done
	"VaultUpdateMinUserDepositEvent":        mapVaultUpdateMinUserDepositEvent,        //Done
	"VaultUpdateProfitMaxUnlockTimeEvent":   mapVaultUpdateProfitMaxUnlockTimeEvent,   //Done
	"VaultUpdateUserDepositLimitEvent":      mapVaultUpdateUserDepositLimitEvent,      //Done
	"VaultUpdateWhitelistedOnlyEvent":       mapVaultUpdateWhitelistedOnlyEvent,       //Done
	"VaultWithdrawlEvent":                   mapVaultWithdrawlEvent,                   //Done
	"WhitelistUpdatedEvent":                 mapWhitelistUpdatedEvent,
	"WithdrawalRequestCanceledEvent":        mapWithdrawalRequestCanceledEvent,  //Done
	"WithdrawalRequestFulfilledEvent":       mapWithdrawalRequestFulfilledEvent, //Done
	"WithdrawalRequestedEvent":              mapWithdrawalRequestedEvent,        //Done
}

func mapEvents(ctx context.Context, db *gorm.DB, event core.Event) error {
	if handler, ok := registry[event.Name]; ok {
		return handler(ctx, db, event)
	}
	log.Warnf("No mapping implemented for event: %s", event.Name)
	return nil
}
