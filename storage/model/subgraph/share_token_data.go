package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShareTokenData struct {
	ID         string       `gorm:"primaryKey;column:id"`         // ID (Int8)
	Vault      *Vault       `gorm:"foreignKey:VaultID"`           // Reference to Vault
	VaultID    string       `gorm:"column:vault_id"`              // Vault ID
	Timestamp  types.BigInt `gorm:"column:timestamp;default:0"`   // Timestamp of the record
	SharePrice types.BigInt `gorm:"column:share_price;default:0"` // Share price (BigInt)
}

func (ShareTokenData) TableName() string {
	return "share_token_data"
}

func (s *ShareTokenData) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		First(s, "id = ?", s.ID).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (s *ShareTokenData) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(s).Error
}
