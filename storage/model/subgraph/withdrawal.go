package subgraph

type Withdrawal struct {
	ID           string   `gorm:"primaryKey;column:id"`          // Transaction-Log
	Timestamp    string   `gorm:"column:timestamp;default:0"`    // Timestamp of update (BigInt)
	BlockNumber  string   `gorm:"column:block_number;default:0"` // Block number of update (BigInt)
	Account      *Account `gorm:"foreignKey:AccountID"`          // Account making withdraw
	AccountID    string   `gorm:"column:account_id"`             // Account ID
	Vault        *Vault   `gorm:"foreignKey:VaultID"`            // Vault withdrawn from
	VaultID      string   `gorm:"column:vault_id"`               // Vault ID
	TokenAmount  string   `gorm:"column:token_amount;default:0"` // Number of Tokens withdrawn from Vault (BigInt)
	SharesBurnt  string   `gorm:"column:shares_burnt;default:0"` // Number of Vault Shares burnt (BigInt)
	Token        *Token   `gorm:"foreignKey:TokenID"`            // Token this Vault will accrue
	TokenID      string   `gorm:"column:token_id"`               // Token ID
	ShareToken   *Token   `gorm:"foreignKey:ShareTokenID"`       // Token representing Shares in the Vault
	ShareTokenID string   `gorm:"column:share_token_id"`         // Share Token ID
	SharePrice   string   `gorm:"column:share_price;default:0"`  // Share price (BigInt)
}

func (Withdrawal) TableName() string {
	return "withdrawals"
}
