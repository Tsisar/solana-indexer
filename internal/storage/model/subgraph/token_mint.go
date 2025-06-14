package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type TokenMint struct {
	ID     string           `gorm:"primaryKey;column:id"` // Mint ID
	Mint   *Token           `gorm:"foreignKey:MintID"`    // Mint account (Token)
	MintID string           `gorm:"column:mint_id"`       // Token mint address
	To     *ShareToken      `gorm:"foreignKey:ToID"`      // Mint account (ShareToken)
	ToID   string           `gorm:"column:to_id"`         // ShareToken ID
	Amount types.BigDecimal `gorm:"column:amount"`        // Number of Tokens minted (BigDecimal)
}

func (TokenMint) TableName() string {
	return "token_mints"
}

func (t *TokenMint) Init() {
	t.MintID = ""
	t.To = nil
	t.ToID = ""
	t.Amount.Zero()
}

func (t *TokenMint) GetID() string {
	return t.ID
}

func (t *TokenMint) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, t)
}

func (t *TokenMint) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, t)
}
