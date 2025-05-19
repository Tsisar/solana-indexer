package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Deposit struct {
	ID           string       `gorm:"primaryKey;column:id"`           // Transaction-Log
	Timestamp    types.BigInt `gorm:"column:timestamp;default:0"`     // Timestamp of update (BigInt)
	BlockNumber  types.BigInt `gorm:"column:block_number;default:0"`  // Block number of update (BigInt)
	Account      *Account     `gorm:"foreignKey:AccountID"`           // Account making Deposit
	AccountID    string       `gorm:"column:account_id"`              // Account ID
	Vault        *Vault       `gorm:"foreignKey:VaultID"`             // Vault deposited into
	VaultID      string       `gorm:"column:vault_id"`                // Vault ID
	TokenAmount  types.BigInt `gorm:"column:token_amount;default:0"`  // Number of Tokens deposited into Vault (BigInt)
	SharesMinted types.BigInt `gorm:"column:shares_minted;default:0"` // Number of new Vault Shares minted (BigInt)
	Token        *Token       `gorm:"foreignKey:TokenID"`             // Token this Vault will accrue
	TokenID      string       `gorm:"column:token_id"`                // Token ID
	ShareToken   *Token       `gorm:"foreignKey:ShareTokenID"`        // Token representing Shares in the Vault
	ShareTokenID string       `gorm:"column:share_token_id"`          // Share Token ID
	SharePrice   types.BigInt `gorm:"column:share_price;default:0"`   // Share price (BigInt)
}

func (Deposit) TableName() string {
	return "deposits"
}

func (d *Deposit) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", d.ID).
		First(d).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (d *Deposit) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(d).Error
}
