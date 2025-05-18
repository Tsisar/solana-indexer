package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VaultHistoricalApr struct {
	ID        string           `gorm:"primaryKey;column:id"`       // ID
	Timestamp types.BigInt     `gorm:"column:timestamp;default:0"` // Time in UTC (BigInt)
	Apr       types.BigDecimal `gorm:"column:apr;default:0"`       // The Annual Percentage Rate (BigDecimal)
	Vault     *Vault           `gorm:"foreignKey:VaultID"`         // The Vault
	VaultID   string           `gorm:"column:vault_id"`            // Vault ID
}

func (VaultHistoricalApr) TableName() string {
	return "vault_historical_aprs"
}

func (v *VaultHistoricalApr) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		First(v, "id = ?", v.ID).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (v *VaultHistoricalApr) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(v).Error
}
