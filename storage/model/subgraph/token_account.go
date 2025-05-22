package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"gorm.io/gorm"
)

type TokenAccount struct {
	ID    string `gorm:"primaryKey;column:id"` // Token account ID
	Mint  string `gorm:"column:mint_id"`       // Mint ID
	Owner string `gorm:"column:owner_id"`      // Owner ID
}

func (TokenAccount) TableName() string {
	return "token_accounts"
}

func (t *TokenAccount) Init() {
	t.Mint = ""
	t.Owner = ""
}

func (t *TokenAccount) GetID() string {
	return t.ID
}

func (t *TokenAccount) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, t)
}

func (t *TokenAccount) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, t)
}
