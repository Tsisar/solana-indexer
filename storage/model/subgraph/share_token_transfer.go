package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShareTokenTransfer struct {
	ID     string           `gorm:"primaryKey;column:id"`    // ID
	To     *Token           `gorm:"foreignKey:ToID"`         // Transfer to account (ShareToken)
	ToID   string           `gorm:"column:to_id"`            // ShareToken ID (recipient)
	From   *Token           `gorm:"foreignKey:FromID"`       // Transfer from account (ShareToken)
	FromID string           `gorm:"column:from_id"`          // ShareToken ID (sender)
	Amount types.BigDecimal `gorm:"column:amount;default:0"` // Number of Tokens transferred (BigDecimal)
}

func (ShareTokenTransfer) TableName() string {
	return "share_token_transfers"
}

func (s *ShareTokenTransfer) Init() {
	s.Amount = types.ZeroBigDecimal()
}

func (s *ShareTokenTransfer) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", s.ID).
		First(s).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		s.Init()
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (s *ShareTokenTransfer) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(s).Error
}
