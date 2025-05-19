package subgraph

import (
	"context"
	"errors"

	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountVaultPosition struct {
	ID              string       `gorm:"primaryKey;column:id"`              // Account-Vault ID
	Vault           *Vault       `gorm:"foreignKey:VaultID"`                // Vault
	VaultID         string       `gorm:"column:vault_id"`                   // Vault ID
	Account         *Account     `gorm:"foreignKey:AccountID"`              // Account
	AccountID       string       `gorm:"column:account_id"`                 // Account ID
	Token           *Token       `gorm:"foreignKey:TokenID"`                // Vault token
	TokenID         string       `gorm:"column:token_id"`                   // Token ID
	ShareToken      *Token       `gorm:"foreignKey:ShareTokenID"`           // Vault share token
	ShareTokenID    string       `gorm:"column:share_token_id"`             // Share Token ID
	BalanceShares   types.BigInt `gorm:"column:balance_shares;default:0"`   // Share balance (BigInt)
	BalanceTokens   types.BigInt `gorm:"column:balance_tokens;default:0"`   // Current token balance (BigInt)
	BalancePosition types.BigInt `gorm:"column:balance_position;default:0"` // Computed position balance (BigInt)
	BalanceProfit   types.BigInt `gorm:"column:balance_profit;default:0"`   // Accumulated profit on full withdrawal (BigInt)
}

func (AccountVaultPosition) TableName() string {
	return "account_vault_positions"
}

func (p *AccountVaultPosition) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", p.ID).
		First(p).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (p *AccountVaultPosition) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(p).Error
}
