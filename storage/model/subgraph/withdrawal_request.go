package subgraph

type WithdrawalRequest struct {
	ID           string  `gorm:"primaryKey;column:id"`           // ID
	User         string  `gorm:"column:user"`                    // User requested withdrawal
	Vault        *Vault  `gorm:"foreignKey:VaultID"`             // Vault requested withdrawal
	VaultID      string  `gorm:"column:vault_id"`                // Vault ID
	Index        string  `gorm:"column:index;default:0"`         // Index (BigInt)
	Recipient    string  `gorm:"column:recipient"`               // Recipient
	Shares       string  `gorm:"column:shares;default:0"`        // Shares (BigInt)
	Amount       string  `gorm:"column:amount;default:0"`        // Amount (BigInt)
	MaxLoss      string  `gorm:"column:max_loss;default:0"`      // Max loss (BigInt)
	FeeShares    string  `gorm:"column:fee_shares;default:0"`    // Fee shares (BigInt)
	Open         bool    `gorm:"column:open"`                    // Flag to check if request is open
	Status       string  `gorm:"column:status"`                  // Status: open/cancelled/fulfilled
	Timestamp    string  `gorm:"column:timestamp;default:0"`     // Timestamp of update (BigInt)
	PriorityFees *string `gorm:"column:priority_fees;default:0"` // Priority fees (BigInt, optional)
}

func (WithdrawalRequest) TableName() string {
	return "withdrawal_requests"
}
