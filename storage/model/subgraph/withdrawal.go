package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Withdrawal struct {
	ID           string       `gorm:"primaryKey;column:id"`          // Transaction-Log
	Timestamp    types.BigInt `gorm:"column:timestamp;default:0"`    // Timestamp of update (BigInt)
	BlockNumber  types.BigInt `gorm:"column:block_number;default:0"` // Block number of update (BigInt)
	Account      *Account     `gorm:"foreignKey:AccountID"`          // Account making withdraw
	AccountID    string       `gorm:"column:account_id"`             // Account ID
	Vault        *Vault       `gorm:"foreignKey:VaultID"`            // Vault withdrawn from
	VaultID      string       `gorm:"column:vault_id"`               // Vault ID
	TokenAmount  types.BigInt `gorm:"column:token_amount;default:0"` // Number of Tokens withdrawn from Vault (BigInt)
	SharesBurnt  types.BigInt `gorm:"column:shares_burnt;default:0"` // Number of Vault Shares burnt (BigInt)
	Token        *Token       `gorm:"foreignKey:TokenID"`            // Token this Vault will accrue
	TokenID      string       `gorm:"column:token_id"`               // Token ID
	ShareToken   *Token       `gorm:"foreignKey:ShareTokenID"`       // Token representing Shares in the Vault
	ShareTokenID string       `gorm:"column:share_token_id"`         // Share Token ID
	SharePrice   types.BigInt `gorm:"column:share_price;default:0"`  // Share price (BigInt)
}

func (Withdrawal) TableName() string {
	return "withdrawals"
}

func (w *Withdrawal) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", w.ID).
		First(w).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (w *Withdrawal) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(w).Error
}
