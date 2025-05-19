package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Token struct {
	ID           string       `gorm:"primaryKey;column:id"`           // Token mint address
	Decimals     types.BigInt `gorm:"column:decimals;default:0"`      // Number of decimals
	Name         string       `gorm:"column:name"`                    // Name of the token
	Symbol       string       `gorm:"column:symbol"`                  // Symbol of the token
	CurrentPrice types.BigInt `gorm:"column:current_price;default:0"` // BigInt → string
}

func (Token) TableName() string {
	return "tokens"
}

func (t *Token) Load(ctx context.Context, db *gorm.DB) (bool, error) {
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

func (t *Token) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(t).Error
}
