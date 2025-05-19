package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StrategyHistoricalApr struct {
	ID        string           `gorm:"primaryKey;column:id"` // Unique ID
	Timestamp types.BigInt     `gorm:"column:timestamp"`     // BigInt
	Apr       types.BigDecimal `gorm:"column:apr"`           // BigDecimal

	StrategyID string    `gorm:"column:strategy_id"`    // Foreign key
	Strategy   *Strategy `gorm:"foreignKey:StrategyID"` // Relation to Strategy
}

func (StrategyHistoricalApr) TableName() string {
	return "strategy_historical_aprs"
}

func (s *StrategyHistoricalApr) Load(ctx context.Context, db *gorm.DB) (bool, error) {
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

func (s *StrategyHistoricalApr) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(s).Error
}
