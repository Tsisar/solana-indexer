package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type Token struct {
	ID           string       `gorm:"primaryKey;column:id"` // Token mint address
	Decimals     types.BigInt `gorm:"column:decimals"`      // Number of decimals
	Name         string       `gorm:"column:name"`          // Name of the token
	Symbol       string       `gorm:"column:symbol"`        // Symbol of the token
	CurrentPrice types.BigInt `gorm:"column:current_price"` // BigInt → string
}

func (Token) TableName() string {
	return "tokens"
}

func (t *Token) Init() {
	t.Decimals.Zero()
	t.Name = ""
	t.Symbol = ""
	t.CurrentPrice.Zero()
}

func (t *Token) GetID() string {
	return t.ID
}

func (t *Token) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, t)
}

func (t *Token) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, t)
}
