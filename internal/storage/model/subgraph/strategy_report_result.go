package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type StrategyReportResult struct {
	ID               string           `gorm:"primaryKey;column:id"`        // The Strategy Report Result ID
	Timestamp        types.BigInt     `gorm:"column:timestamp"`            // Timestamp the strategy report was most recently updated (BigInt)
	BlockNumber      types.BigInt     `gorm:"column:block_number"`         // Blocknumber the strategy report was most recently updated (BigInt)
	CurrentReport    *StrategyReport  `gorm:"foreignKey:CurrentReportID"`  // The current strategy report
	CurrentReportID  string           `gorm:"column:current_report_id"`    // Current Report ID
	PreviousReport   *StrategyReport  `gorm:"foreignKey:PreviousReportID"` // The previous strategy report
	PreviousReportID string           `gorm:"column:previous_report_id"`   // Previous Report ID
	StartTimestamp   types.BigInt     `gorm:"column:start_timestamp"`      // Start timestamp (BigInt)
	EndTimestamp     types.BigInt     `gorm:"column:end_timestamp"`        // End timestamp (BigInt)
	Duration         types.BigDecimal `gorm:"column:duration"`             // The duration (in days) from the previous report (BigDecimal)
	DurationPr       types.BigDecimal `gorm:"column:duration_pr"`          // Duration percentage rate (BigDecimal)
	Apr              types.BigDecimal `gorm:"column:apr"`                  // Annual Percentage Rate (BigDecimal)
	TransactionHash  string           `gorm:"column:transaction_hash"`     // Transaction Hash
}

func (StrategyReportResult) TableName() string {
	return "strategy_report_results"
}

func (s *StrategyReportResult) Init() {
	s.Timestamp.Zero()
	s.BlockNumber.Zero()
	s.CurrentReport = nil
	s.CurrentReportID = ""
	s.PreviousReport = nil
	s.PreviousReportID = ""
	s.StartTimestamp.Zero()
	s.EndTimestamp.Zero()
	s.Duration.Zero()
	s.DurationPr.Zero()
	s.Apr.Zero()
	s.TransactionHash = ""
}

func (s *StrategyReportResult) GetID() string {
	return s.ID
}

func (s *StrategyReportResult) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, s)
}

func (s *StrategyReportResult) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, s)
}
