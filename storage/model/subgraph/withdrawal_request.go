package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

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

func (w *WithdrawalRequest) Init() {
	w.User = ""
	w.VaultID = ""
	w.Recipient = ""
	w.Status = ""
}

func (w *WithdrawalRequest) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", w.ID).
		First(w).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		w.Init()
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (w *WithdrawalRequest) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(w).Error
}
