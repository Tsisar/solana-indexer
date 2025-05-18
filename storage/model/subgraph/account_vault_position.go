package subgraph

type AccountVaultPosition struct {
	ID              string   `gorm:"primaryKey;column:id"`              // Account-Vault ID
	Vault           *Vault   `gorm:"foreignKey:VaultID"`                // Vault
	VaultID         string   `gorm:"column:vault_id"`                   // Vault ID
	Account         *Account `gorm:"foreignKey:AccountID"`              // Account
	AccountID       string   `gorm:"column:account_id"`                 // Account ID
	Token           *Token   `gorm:"foreignKey:TokenID"`                // Vault token
	TokenID         string   `gorm:"column:token_id"`                   // Token ID
	ShareToken      *Token   `gorm:"foreignKey:ShareTokenID"`           // Vault share token
	ShareTokenID    string   `gorm:"column:share_token_id"`             // Share Token ID
	BalanceShares   string   `gorm:"column:balance_shares;default:0"`   // Share balance (BigInt)
	BalanceTokens   string   `gorm:"column:balance_tokens;default:0"`   // Current token balance (BigInt)
	BalancePosition string   `gorm:"column:balance_position;default:0"` // Computed position balance (BigInt)
	BalanceProfit   string   `gorm:"column:balance_profit;default:0"`   // Accumulated profit on full withdrawal (BigInt)
}

func (AccountVaultPosition) TableName() string {
	return "account_vault_positions"
}
