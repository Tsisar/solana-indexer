package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DeployFunds struct {
	ID         string       `gorm:"primaryKey;column:id"`       // The Strategy Deploy ID
	Strategy   *Strategy    `gorm:"foreignKey:StrategyID"`      // The Strategy
	StrategyID string       `gorm:"column:strategy_id"`         // Strategy ID
	Amount     types.BigInt `gorm:"column:amount;default:0"`    // Total amount of assets deposited in strategy (BigInt)
	Timestamp  types.BigInt `gorm:"column:timestamp;default:0"` // Timestamp the strategy report was most recently updated (BigInt)
}

func (DeployFunds) TableName() string {
	return "deploy_funds"
}

func (d *DeployFunds) Init() {
	d.Amount = types.ZeroBigInt()
	d.Timestamp = types.ZeroBigInt()
	d.StrategyID = ""
}

func (d *DeployFunds) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", d.ID).
		First(d).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		d.Init()
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (d *DeployFunds) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(d).Error
}
