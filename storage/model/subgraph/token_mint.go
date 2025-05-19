package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenMint struct {
	ID     string           `gorm:"primaryKey;column:id"`    // Mint ID
	To     *Token           `gorm:"foreignKey:ToID"`         // Mint account (ShareToken)
	ToID   string           `gorm:"column:to_id"`            // ShareToken ID
	Amount types.BigDecimal `gorm:"column:amount;default:0"` // Number of Tokens minted (BigDecimal)
}

func (TokenMint) TableName() string {
	return "token_mints"
}

func (t *TokenMint) Init() {
	t.Amount = types.ZeroBigDecimal()
}

func (t *TokenMint) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", t.ID).
		First(t).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		t.Init()
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (t *TokenMint) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(t).Error
}
