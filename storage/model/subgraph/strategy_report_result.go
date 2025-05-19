package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StrategyReportResult struct {
	ID               string           `gorm:"primaryKey;column:id"`             // The Strategy Report Result ID
	Timestamp        types.BigInt     `gorm:"column:timestamp;default:0"`       // Timestamp the strategy report was most recently updated (BigInt)
	BlockNumber      types.BigInt     `gorm:"column:block_number;default:0"`    // Blocknumber the strategy report was most recently updated (BigInt)
	CurrentReport    *StrategyReport  `gorm:"foreignKey:CurrentReportID"`       // The current strategy report
	CurrentReportID  string           `gorm:"column:current_report_id"`         // Current Report ID
	PreviousReport   *StrategyReport  `gorm:"foreignKey:PreviousReportID"`      // The previous strategy report
	PreviousReportID string           `gorm:"column:previous_report_id"`        // Previous Report ID
	StartTimestamp   types.BigInt     `gorm:"column:start_timestamp;default:0"` // Start timestamp (BigInt)
	EndTimestamp     types.BigInt     `gorm:"column:end_timestamp;default:0"`   // End timestamp (BigInt)
	Duration         types.BigDecimal `gorm:"column:duration;default:0"`        // The duration (in days) from the previous report (BigDecimal)
	DurationPr       types.BigDecimal `gorm:"column:duration_pr;default:0"`     // Duration percentage rate (BigDecimal)
	Apr              types.BigDecimal `gorm:"column:apr;default:0"`             // Annual Percentage Rate (BigDecimal)
	TransactionHash  string           `gorm:"column:transaction_hash"`          // Transaction Hash
}

func (StrategyReportResult) TableName() string {
	return "strategy_report_results"
}

func (s *StrategyReportResult) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", s.ID).
		First(s).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (s *StrategyReportResult) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(s).Error
}
