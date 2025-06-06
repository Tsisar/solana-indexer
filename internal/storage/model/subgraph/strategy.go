package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type Strategy struct {
	ID                      string                   `gorm:"primaryKey;column:id"`               // Strategy address
	Vault                   *Vault                   `gorm:"foreignKey:VaultID"`                 // The Vault
	VaultID                 string                   `gorm:"column:vault_id"`                    // Vault ID
	StrategyType            string                   `gorm:"column:strategy_type"`               // Strategy type
	Amount                  types.BigInt             `gorm:"column:amount"`                      // Token this Strategy will accrue (BigInt)
	TotalAssets             types.BigInt             `gorm:"column:total_assets"`                // Total amount of assets deposited in strategies (BigInt)
	TotalInvested           types.BigInt             `gorm:"column:total_invested"`              // Total invested amount (BigInt)
	DepositLimit            types.BigInt             `gorm:"column:deposit_limit"`               // The maximum amount of tokens (BigInt)
	DepositPeriodEnds       types.BigInt             `gorm:"column:deposit_period_ends"`         // Deposit period ends (BigInt)
	LockPeriodEnds          types.BigInt             `gorm:"column:lock_period_ends"`            // Lock period ends (BigInt)
	CurrentDebt             types.BigInt             `gorm:"column:current_debt"`                // Current debt of the strategy (BigInt)
	MaxDebt                 types.BigInt             `gorm:"column:max_debt"`                    // Maximum allowed debt (BigInt)
	Apr                     types.BigDecimal         `gorm:"column:apr"`                         // Annual Percentage Rate (BigDecimal)
	Activation              types.BigInt             `gorm:"column:activation"`                  // Creation timestamp (BigInt)
	DelegatedAssets         *types.BigInt            `gorm:"column:delegated_assets"`            // Delegated assets (BigInt, optional)
	LatestReportID          *string                  `gorm:"column:latest_report_id"`            // Latest Report ID
	Reports                 []*StrategyReport        `gorm:"foreignKey:StrategyID"`              // Reports created by this strategy
	ReportsEvents           []*StrategyReportEvent   `gorm:"foreignKey:StrategyID"`              // Report events created by this strategy
	ReportsCount            types.BigInt             `gorm:"column:reports_count"`               // Reports count (BigDecimal? or BigInt?)
	HistoricalApr           []*StrategyHistoricalApr `gorm:"foreignKey:StrategyID"`              // Historical APR
	PerformanceFees         types.BigInt             `gorm:"column:performance_fees"`            // Protocol fees (BigInt)
	DtfReport               *DTFReport               `gorm:"foreignKey:DtfReportID"`             // DTF Report
	DtfReportID             *string                  `gorm:"column:dtf_report_id"`               // DTF Report ID
	TotalAllocation         types.BigDecimal         `gorm:"column:total_allocation"`            // Total allocation after swap (BigDecimal)
	TotalAllocationPercent  types.BigDecimal         `gorm:"column:total_allocation_in_precent"` // Total allocation percent after swap (BigDecimal)
	EffectiveInvestedAmount types.BigInt             `gorm:"column:effective_invested_amount"`   // Effective invested amount (BigInt)
	ProfitOrLoss            types.BigDecimal         `gorm:"column:profit_or_loss"`              // Profit or loss (BigDecimal)
	ProfitOrLossPercent     types.BigDecimal         `gorm:"column:profit_or_loss_in_precent"`   // Profit/loss in percent (BigDecimal)
	DeployFunds             []*DeployFunds           `gorm:"foreignKey:StrategyID"`              // Deploy Funds
	FreeFunds               []*FreeFunds             `gorm:"foreignKey:StrategyID"`              // Free Funds
	Asset                   *Token                   `gorm:"foreignKey:AssetID"`                 // Asset the strategy is investing in
	AssetID                 *string                  `gorm:"column:asset_id"`                    // Asset ID (optional)
	Removed                 bool                     `gorm:"column:removed"`                     // Removed from vault
	RemovedTimestamp        types.BigInt             `gorm:"column:removed_timestamp"`           // Removed timestamp (BigInt)
	UnderlyingMint          string                   `gorm:"column:underlying_mint"`             // Underlying token mint address
	UnderlyingTokenAcc      string                   `gorm:"column:underlying_token_acc"`        // Underlying token account address
	UnderlyingDecimals      types.BigInt             `gorm:"column:underlying_decimals"`         // Underlying token decimals
}

func (Strategy) TableName() string {
	return "strategies"
}

func (s *Strategy) Init() {
	s.VaultID = ""
	s.StrategyType = ""

	s.Amount.Zero()
	s.TotalAssets.Zero()
	s.TotalInvested.Zero()
	s.DepositLimit.Zero()
	s.DepositPeriodEnds.Zero()
	s.LockPeriodEnds.Zero()
	s.CurrentDebt.Zero()
	s.MaxDebt.Zero()
	s.Apr.Zero()
	s.Activation.Zero()
	s.ReportsCount.Zero()
	s.PerformanceFees.Zero()
	s.TotalAllocation.Zero()
	s.TotalAllocationPercent.Zero()
	s.EffectiveInvestedAmount.Zero()
	s.ProfitOrLoss.Zero()
	s.ProfitOrLossPercent.Zero()
	s.Removed = false
	s.RemovedTimestamp.Zero()
	s.LatestReportID = nil
}

func (s *Strategy) GetID() string {
	return s.ID
}

func (s *Strategy) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, s)
}

func (s *Strategy) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, s)
}
