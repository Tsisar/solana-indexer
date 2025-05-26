package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"gorm.io/gorm"
)

type TokenWallet struct {
	ID          string   `gorm:"primaryKey;column:id"`   // Account address (Associated Token Account)
	Authority   *Account `gorm:"foreignKey:AuthorityID"` // Authority
	AuthorityID string   `gorm:"column:authority_id"`    // Authority ID
}

func (TokenWallet) TableName() string {
	return "token_wallets"
}

func (t *TokenWallet) Init() {
	t.Authority = nil
	t.AuthorityID = ""
}

func (t *TokenWallet) GetID() string {
	return t.ID
}

func (t *TokenWallet) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, t)
}

func (t *TokenWallet) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, t)
}
