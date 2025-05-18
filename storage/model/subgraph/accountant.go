package subgraph

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Accountant struct {
	ID              string `gorm:"primaryKey;column:id"`              // Accountant address
	EntryFee        string `gorm:"column:entry_fee;default:0"`        // Entry fee
	RedemptionFee   string `gorm:"column:redemption_fee;default:0"`   // Redemption fee
	PerformanceFees string `gorm:"column:performance_fees;default:0"` // Performance fees
}

func (*Accountant) TableName() string {
	return "accountants"
}

func (a *Accountant) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", a.ID).
		First(a).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (a *Accountant) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(a).Error
}
