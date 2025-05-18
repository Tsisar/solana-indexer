package subgraph

import "github.com/Tsisar/solana-indexer/subgraph/types"

type WithdrawalRequest struct {
	ID           string        `gorm:"primaryKey;column:id"`           // ID
	User         string        `gorm:"column:user"`                    // User requested withdrawal
	Vault        *Vault        `gorm:"foreignKey:VaultID"`             // Vault requested withdrawal
	VaultID      string        `gorm:"column:vault_id"`                // Vault ID
	Index        types.BigInt  `gorm:"column:index;default:0"`         // Index (BigInt)
	Recipient    string        `gorm:"column:recipient"`               // Recipient
	Shares       types.BigInt  `gorm:"column:shares;default:0"`        // Shares (BigInt)
	Amount       types.BigInt  `gorm:"column:amount;default:0"`        // Amount (BigInt)
	MaxLoss      types.BigInt  `gorm:"column:max_loss;default:0"`      // Max loss (BigInt)
	FeeShares    types.BigInt  `gorm:"column:fee_shares;default:0"`    // Fee shares (BigInt)
	Open         bool          `gorm:"column:open"`                    // Flag to check if request is open
	Status       string        `gorm:"column:status"`                  // Status: open/cancelled/fulfilled
	Timestamp    types.BigInt  `gorm:"column:timestamp;default:0"`     // Timestamp of update (BigInt)
	PriorityFees *types.BigInt `gorm:"column:priority_fees;default:0"` // Priority fees (BigInt, optional)
}

func (WithdrawalRequest) TableName() string {
	return "withdrawal_requests"
}
