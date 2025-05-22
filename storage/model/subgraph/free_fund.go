package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type FreeFunds struct {
	ID         string       `gorm:"primaryKey;column:id"`  // The Strategy Fund ID
	Strategy   *Strategy    `gorm:"foreignKey:StrategyID"` // The Strategy
	StrategyID string       `gorm:"column:strategy_id"`    // Strategy ID
	Amount     types.BigInt `gorm:"column:amount"`         // Total amount of assets deposited in strategies (BigInt)
	Timestamp  types.BigInt `gorm:"column:timestamp"`      // Timestamp the strategy report was most recently updated (BigInt)
}

func (FreeFunds) TableName() string {
	return "free_funds"
}

func (f *FreeFunds) Init() {
	f.Strategy = nil
	f.StrategyID = ""
	f.Amount.Zero()
	f.Timestamp.Zero()
}

func (f *FreeFunds) GetID() string {
	return f.ID
}

func (f *FreeFunds) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, f)
}

func (f *FreeFunds) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, f)
}
