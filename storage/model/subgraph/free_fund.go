package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FreeFunds struct {
	ID         string       `gorm:"primaryKey;column:id"`       // The Strategy Fund ID
	Strategy   *Strategy    `gorm:"foreignKey:StrategyID"`      // The Strategy
	StrategyID string       `gorm:"column:strategy_id"`         // Strategy ID
	Amount     types.BigInt `gorm:"column:amount;default:0"`    // Total amount of assets deposited in strategies (BigInt)
	Timestamp  types.BigInt `gorm:"column:timestamp;default:0"` // Timestamp the strategy report was most recently updated (BigInt)
}

func (FreeFunds) TableName() string {
	return "free_funds"
}

func (f *FreeFunds) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", f.ID).
		First(f).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (f *FreeFunds) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(f).Error
}
