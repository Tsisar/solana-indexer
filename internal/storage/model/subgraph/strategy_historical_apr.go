package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	types2 "github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type StrategyHistoricalApr struct {
	ID        string            `gorm:"primaryKey;column:id"` // Unique ID
	Timestamp types2.BigInt     `gorm:"column:timestamp"`     // BigInt
	Apr       types2.BigDecimal `gorm:"column:apr"`           // BigDecimal

	StrategyID string    `gorm:"column:strategy_id"`    // Foreign key
	Strategy   *Strategy `gorm:"foreignKey:StrategyID"` // Relation to Strategy
}

func (StrategyHistoricalApr) TableName() string {
	return "strategy_historical_aprs"
}

func (s *StrategyHistoricalApr) Init() {
	s.Timestamp.Zero()
	s.Apr.Zero()
	s.StrategyID = ""
	s.Strategy = nil
}

func (s *StrategyHistoricalApr) GetID() string {
	return s.ID
}

func (s *StrategyHistoricalApr) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, s)
}

func (s *StrategyHistoricalApr) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, s)
}
