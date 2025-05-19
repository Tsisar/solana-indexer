package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DTFReport struct {
	ID          string       `gorm:"primaryKey;column:id"`          // The Strategy Report ID
	TotalAssets types.BigInt `gorm:"column:total_assets;default:0"` // Total amount of assets deposited in strategies (BigInt)
	Timestamp   types.BigInt `gorm:"column:timestamp;default:0"`    // Timestamp the strategy report was most recently updated (BigInt)
}

func (DTFReport) TableName() string {
	return "dtf_reports"
}

func (d *DTFReport) Init() {
	d.TotalAssets = types.ZeroBigInt()
	d.Timestamp = types.ZeroBigInt()
}

func (d *DTFReport) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", d.ID).
		First(d).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		d.Init()
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (d *DTFReport) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(d).Error
}
