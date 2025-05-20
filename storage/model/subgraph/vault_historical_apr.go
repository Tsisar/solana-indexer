package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type VaultHistoricalApr struct {
	ID        string           `gorm:"primaryKey;column:id"` // ID
	Timestamp types.BigInt     `gorm:"column:timestamp"`     // Time in UTC (BigInt)
	Apr       types.BigDecimal `gorm:"column:apr"`           // The Annual Percentage Rate (BigDecimal)
	Vault     *Vault           `gorm:"foreignKey:VaultID"`   // The Vault
	VaultID   string           `gorm:"column:vault_id"`      // Vault ID
}

func (VaultHistoricalApr) TableName() string {
	return "vault_historical_aprs"
}

func (v *VaultHistoricalApr) Init() {
	v.Timestamp.Zero()
	v.Apr.Zero()
	v.Vault = nil
	v.VaultID = ""
}

func (v *VaultHistoricalApr) GetID() string {
	return v.ID
}

func (v *VaultHistoricalApr) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, v)
}

func (v *VaultHistoricalApr) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, v)
}
