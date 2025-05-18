package subgraph

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Strategy struct {
	ID                string  `gorm:"primaryKey;column:id"`                 // Strategy address
	Vault             *Vault  `gorm:"foreignKey:VaultID"`                   // The Vault
	VaultID           string  `gorm:"column:vault_id"`                      // Vault ID
	StrategyType      string  `gorm:"column:strategy_type"`                 // Strategy type
	Amount            string  `gorm:"column:amount;default:0"`              // Token this Strategy will accrue (BigInt)
	TotalAssets       string  `gorm:"column:total_assets;default:0"`        // Total amount of assets deposited in strategies (BigInt)
	DepositLimit      string  `gorm:"column:deposit_limit;default:0"`       // The maximum amount of tokens (BigInt)
	DepositPeriodEnds string  `gorm:"column:deposit_period_ends;default:0"` // Deposit period ends (BigInt)
	LockPeriodEnds    string  `gorm:"column:lock_period_ends;default:0"`    // Lock period ends (BigInt)
	CurrentDebt       string  `gorm:"column:current_debt;default:0"`        // Current debt of the strategy (BigInt)
	MaxDebt           string  `gorm:"column:max_debt;default:0"`            // Maximum allowed debt (BigInt)
	Apr               string  `gorm:"column:apr;default:0"`                 // Annual Percentage Rate (BigDecimal)
	Activation        string  `gorm:"column:activation;default:0"`          // Creation timestamp (BigInt)
	DelegatedAssets   *string `gorm:"column:delegated_assets;default:0"`    // Delegated assets (BigInt, optional)
	//LatestReport            *StrategyReport          `gorm:"foreignKey:LatestReportID"`                    // The latest report for this Strategy
	LatestReportID          *string                  `gorm:"column:latest_report_id"`                      // Latest Report ID
	Reports                 []*StrategyReport        `gorm:"foreignKey:StrategyID"`                        // Reports created by this strategy
	ReportsEvents           []*StrategyReportEvent   `gorm:"foreignKey:StrategyID"`                        // Report events created by this strategy
	ReportsCount            string                   `gorm:"column:reports_count;default:0"`               // Reports count (BigDecimal)
	HistoricalApr           []*StrategyHistoricalApr `gorm:"foreignKey:StrategyID"`                        // Historical APR
	PerformanceFees         string                   `gorm:"column:performance_fees;default:0"`            // Protocol fees (BigInt)
	DtfReport               *DTFReport               `gorm:"foreignKey:DtfReportID"`                       // DTF Report
	DtfReportID             *string                  `gorm:"column:dtf_report_id"`                         // DTF Report ID
	TotalAllocation         string                   `gorm:"column:total_allocation;default:0"`            // Total allocation after swap (BigDecimal)
	TotalAllocationPercent  string                   `gorm:"column:total_allocation_in_precent;default:0"` // Total allocation percent after swap (BigDecimal)
	EffectiveInvestedAmount string                   `gorm:"column:effective_invested_amount;default:0"`   // Effective invested amount (BigInt)
	ProfitOrLoss            string                   `gorm:"column:profit_or_loss;default:0"`              // Profit or loss (BigDecimal)
	ProfitOrLossPercent     string                   `gorm:"column:profit_or_loss_in_precent;default:0"`   // Profit/loss in percent (BigDecimal)
	DeployFunds             []*DeployFunds           `gorm:"foreignKey:StrategyID"`                        // Deploy Funds
	FreeFunds               []*FreeFunds             `gorm:"foreignKey:StrategyID"`                        // Free Funds
	Asset                   *Token                   `gorm:"foreignKey:AssetID"`                           // Asset the strategy is investing in
	AssetID                 *string                  `gorm:"column:asset_id"`                              // Asset ID (optional)
	Removed                 bool                     `gorm:"column:removed"`                               // Removed from vault
	RemovedTimestamp        string                   `gorm:"column:removed_timestamp;default:0"`           // Removed timestamp (BigInt)
}

func (*Strategy) TableName() string {
	return "strategies"
}

func (s *Strategy) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		First(s, "id = ?", s.ID).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (s *Strategy) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(s).Error
}
