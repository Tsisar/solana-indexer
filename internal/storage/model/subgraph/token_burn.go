package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type TokenBurn struct {
	ID     string           `gorm:"primaryKey;column:id"` // Burn ID
	Mint   *Token           `gorm:"foreignKey:MintID"`    // Mint account (Token)
	MintID string           `gorm:"column:mint_id"`       // Token mint address
	From   *ShareToken      `gorm:"foreignKey:FromID"`    // Burn account (ShareToken)
	FromID string           `gorm:"column:from_id"`       // ShareToken ID
	Amount types.BigDecimal `gorm:"column:amount"`        // Number of Tokens burnt (BigDecimal)
}

func (TokenBurn) TableName() string {
	return "token_burns"
}

func (t *TokenBurn) Init() {
	t.MintID = ""
	t.From = nil
	t.FromID = ""
	t.Amount.Zero()
}

func (t *TokenBurn) GetID() string {
	return t.ID
}

func (t *TokenBurn) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, t)
}

func (t *TokenBurn) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, t)
}
