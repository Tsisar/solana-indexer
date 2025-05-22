package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type WithdrawalRequest struct {
	ID           string        `gorm:"primaryKey;column:id"` // ID
	User         string        `gorm:"column:user"`          // User requested withdrawal
	Vault        *Vault        `gorm:"foreignKey:VaultID"`   // Vault requested withdrawal
	VaultID      string        `gorm:"column:vault_id"`      // Vault ID
	Index        types.BigInt  `gorm:"column:index"`         // Index (BigInt)
	Recipient    string        `gorm:"column:recipient"`     // Recipient
	Shares       types.BigInt  `gorm:"column:shares"`        // Shares (BigInt)
	Amount       types.BigInt  `gorm:"column:amount"`        // Amount (BigInt)
	MaxLoss      types.BigInt  `gorm:"column:max_loss"`      // Max loss (BigInt)
	FeeShares    types.BigInt  `gorm:"column:fee_shares"`    // Fee shares (BigInt)
	Open         bool          `gorm:"column:open"`          // Flag to check if request is open
	Status       string        `gorm:"column:status"`        // Status: open/cancelled/fulfilled
	Timestamp    types.BigInt  `gorm:"column:timestamp"`     // Timestamp of update (BigInt)
	PriorityFees *types.BigInt `gorm:"column:priority_fees"` // Priority fees (BigInt, optional)
}

func (WithdrawalRequest) TableName() string {
	return "withdrawal_requests"
}

func (w *WithdrawalRequest) Init() {
	w.User = ""
	w.Vault = nil
	w.VaultID = ""
	w.Index.Zero()
	w.Recipient = ""
	w.Shares.Zero()
	w.Amount.Zero()
	w.MaxLoss.Zero()
	w.FeeShares.Zero()
	w.Open = false
	w.Status = ""
	w.Timestamp.Zero()
	if w.PriorityFees != nil {
		w.PriorityFees.Zero()
	}
}

func (w *WithdrawalRequest) GetID() string {
	return w.ID
}

func (w *WithdrawalRequest) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, w)
}

func (w *WithdrawalRequest) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, w)
}
