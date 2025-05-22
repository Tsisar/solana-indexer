package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type TokenStats struct {
	ID         string           `gorm:"primaryKey;column:id"`  // Aggregated ID (Int8)
	Timestamp  types.BigInt     `gorm:"column:timestamp"`      // Timestamp of aggregation
	Vault      *Vault           `gorm:"foreignKey:VaultID"`    // Vault reference
	VaultID    string           `gorm:"column:vault_id"`       // Vault ID
	SharePrice types.BigDecimal `gorm:"column:share_price"`    // Aggregated share price (last) (BigDecimal)
	Interval   string           `gorm:"not null;default:hour"` // "hour", "day" etc.
}

func (TokenStats) TableName() string {
	return "token_stats"
}

func (t *TokenStats) Init() {
	t.Timestamp.Zero()
	t.Vault = nil
	t.VaultID = ""
	t.SharePrice.Zero()
}

func (t *TokenStats) GetID() string {
	return t.ID
}

func (t *TokenStats) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, t)
}

func (t *TokenStats) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, t)
}
