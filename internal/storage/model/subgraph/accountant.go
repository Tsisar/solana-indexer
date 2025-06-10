package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type Accountant struct {
	ID              string       `gorm:"primaryKey;column:id"`    // Accountant address
	EntryFee        types.BigInt `gorm:"column:entry_fee"`        // Entry fee
	RedemptionFee   types.BigInt `gorm:"column:redemption_fee"`   // Redemption fee
	PerformanceFees types.BigInt `gorm:"column:performance_fees"` // Performance fees
}

func (Accountant) TableName() string {
	return "accountants"
}

func (a *Accountant) Init() {
	a.EntryFee.Zero()
	a.RedemptionFee.Zero()
	a.PerformanceFees.Zero()
}

func (a *Accountant) GetID() string {
	return a.ID
}

func (a *Accountant) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, a)
}

func (a *Accountant) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, a)
}
