package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type StrategyReportEvent struct {
	ID              string       `gorm:"primaryKey;column:id"`    // The Strategy Report ID
	Timestamp       types.BigInt `gorm:"column:timestamp"`        // Timestamp the strategy report was most recently updated (BigInt)
	BlockNumber     types.BigInt `gorm:"column:block_number"`     // Blocknumber the strategy report was most recently updated (BigInt)
	TransactionHash string       `gorm:"column:transaction_hash"` // Transaction Hash
	Strategy        *Strategy    `gorm:"foreignKey:StrategyID"`   // The Strategy reference
	StrategyID      string       `gorm:"column:strategy_id"`      // Strategy ID
	SharePrice      types.BigInt `gorm:"column:share_price"`      // Share price (BigInt)
}

func (StrategyReportEvent) TableName() string {
	return "strategy_report_events"
}

func (s *StrategyReportEvent) Init() {
	s.Timestamp.Zero()
	s.BlockNumber.Zero()
	s.TransactionHash = ""
	s.Strategy = nil
	s.StrategyID = ""
	s.SharePrice.Zero()
}

func (s *StrategyReportEvent) GetID() string {
	return s.ID
}

func (s *StrategyReportEvent) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, s)
}

func (s *StrategyReportEvent) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, s)
}
