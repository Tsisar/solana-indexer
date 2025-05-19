package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenStats struct {
	ID         string           `gorm:"primaryKey;column:id"`         // Aggregated ID (Int8)
	Timestamp  string           `gorm:"column:timestamp;default:0"`   // Timestamp of aggregation
	Vault      *Vault           `gorm:"foreignKey:VaultID"`           // Vault reference
	VaultID    string           `gorm:"column:vault_id"`              // Vault ID
	SharePrice types.BigDecimal `gorm:"column:share_price;default:0"` // Aggregated share price (last) (BigDecimal)
}

func (TokenStats) TableName() string {
	return "token_stats"
}

func (t *TokenStats) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", t.ID).
		First(t).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (t *TokenStats) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(t).Error
}
