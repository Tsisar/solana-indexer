package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type DTFReport struct {
	ID          string       `gorm:"primaryKey;column:id"` // The Strategy Report ID
	TotalAssets types.BigInt `gorm:"column:total_assets"`  // Total amount of assets deposited in strategies (BigInt)
	Timestamp   types.BigInt `gorm:"column:timestamp"`     // Timestamp the strategy report was most recently updated (BigInt)
}

func (DTFReport) TableName() string {
	return "dtf_reports"
}

func (d *DTFReport) Init() {
	d.TotalAssets.Zero()
	d.Timestamp.Zero()
}

func (d *DTFReport) GetID() string {
	return d.ID
}

func (d *DTFReport) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, d)
}

func (d *DTFReport) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, d)
}
