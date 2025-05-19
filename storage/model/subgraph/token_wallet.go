package subgraph

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenWallet struct {
	ID          string   `gorm:"primaryKey;column:id"`   // Account address (Associated Token Account)
	Authority   *Account `gorm:"foreignKey:AuthorityID"` // Authority
	AuthorityID string   `gorm:"column:authority_id"`    // Authority ID
}

func (TokenWallet) TableName() string {
	return "token_wallets"
}

func (t *TokenWallet) Load(ctx context.Context, db *gorm.DB) (bool, error) {
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

func (t *TokenWallet) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(t).Error
}
