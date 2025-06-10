package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type DeployFunds struct {
	ID         string       `gorm:"primaryKey;column:id"`  // The Strategy Deploy ID
	Strategy   *Strategy    `gorm:"foreignKey:StrategyID"` // The Strategy
	StrategyID string       `gorm:"column:strategy_id"`    // Strategy ID
	Amount     types.BigInt `gorm:"column:amount"`         // Total amount of assets deposited in strategy (BigInt)
	Timestamp  types.BigInt `gorm:"column:timestamp"`      // Timestamp the strategy report was most recently updated (BigInt)
}

func (DeployFunds) TableName() string {
	return "deploy_funds"
}

func (d *DeployFunds) Init() {
	d.Strategy = nil
	d.StrategyID = ""
	d.Amount.Zero()
	d.Timestamp.Zero()
}

func (d *DeployFunds) GetID() string {
	return d.ID
}

func (d *DeployFunds) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, d)
}

func (d *DeployFunds) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, d)
}
