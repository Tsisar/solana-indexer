package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Accountant struct {
	ID              string       `gorm:"primaryKey;column:id"`              // Accountant address
	EntryFee        types.BigInt `gorm:"column:entry_fee;default:0"`        // Entry fee
	RedemptionFee   types.BigInt `gorm:"column:redemption_fee;default:0"`   // Redemption fee
	PerformanceFees types.BigInt `gorm:"column:performance_fees;default:0"` // Performance fees
}

func (Accountant) TableName() string {
	return "accountants"
}

func (a *Accountant) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", a.ID).
		First(a).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		a.Init()
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (a *Accountant) Init() {
	a.EntryFee = types.ZeroBigInt()
	a.RedemptionFee = types.ZeroBigInt()
	a.PerformanceFees = types.ZeroBigInt()
}

func (a *Accountant) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(a).Error
}
