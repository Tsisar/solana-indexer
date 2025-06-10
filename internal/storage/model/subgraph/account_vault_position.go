package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type AccountVaultPosition struct {
	ID              string       `gorm:"primaryKey;column:id"`    // Account-Vault ID
	Vault           *Vault       `gorm:"foreignKey:VaultID"`      // Vault
	VaultID         string       `gorm:"column:vault_id"`         // Vault ID
	Account         *Account     `gorm:"foreignKey:AccountID"`    // Account
	AccountID       string       `gorm:"column:account_id"`       // Account ID
	Token           *Token       `gorm:"foreignKey:TokenID"`      // Vault token
	TokenID         string       `gorm:"column:token_id"`         // Token ID
	ShareToken      *Token       `gorm:"foreignKey:ShareTokenID"` // Vault share token
	ShareTokenID    string       `gorm:"column:share_token_id"`   // Share Token ID
	BalanceShares   types.BigInt `gorm:"column:balance_shares"`   // Share balance (BigInt)
	BalanceTokens   types.BigInt `gorm:"column:balance_tokens"`   // Current token balance (BigInt)
	BalancePosition types.BigInt `gorm:"column:balance_position"` // Computed position balance (BigInt)
	BalanceProfit   types.BigInt `gorm:"column:balance_profit"`   // Accumulated profit on full withdrawal (BigInt)
}

func (AccountVaultPosition) TableName() string {
	return "account_vault_positions"
}

func (p *AccountVaultPosition) Init() {
	p.VaultID = ""
	p.AccountID = ""
	p.TokenID = ""
	p.ShareTokenID = ""

	p.BalanceShares.Zero()
	p.BalanceTokens.Zero()
	p.BalancePosition.Zero()
	p.BalanceProfit.Zero()
}

func (p *AccountVaultPosition) GetID() string {
	return p.ID
}

func (p *AccountVaultPosition) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, p)
}

func (p *AccountVaultPosition) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, p)
}
