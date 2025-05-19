package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenBurn struct {
	ID     string           `gorm:"primaryKey;column:id"`    // Burn ID
	From   *Token           `gorm:"foreignKey:FromID"`       // Burn account (ShareToken)
	FromID string           `gorm:"column:from_id"`          // ShareToken ID
	Amount types.BigDecimal `gorm:"column:amount;default:0"` // Number of Tokens burnt (BigDecimal)
}

func (TokenBurn) TableName() string {
	return "token_burns"
}

func (t *TokenBurn) Load(ctx context.Context, db *gorm.DB) (bool, error) {
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

func (t *TokenBurn) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(t).Error
}
