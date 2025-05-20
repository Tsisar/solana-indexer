package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type Withdrawal struct {
	ID           string       `gorm:"primaryKey;column:id"`    // Transaction-Log
	Timestamp    types.BigInt `gorm:"column:timestamp"`        // Timestamp of update (BigInt)
	BlockNumber  types.BigInt `gorm:"column:block_number"`     // Block number of update (BigInt)
	Account      *Account     `gorm:"foreignKey:AccountID"`    // Account making withdraw
	AccountID    string       `gorm:"column:account_id"`       // Account ID
	Vault        *Vault       `gorm:"foreignKey:VaultID"`      // Vault withdrawn from
	VaultID      string       `gorm:"column:vault_id"`         // Vault ID
	TokenAmount  types.BigInt `gorm:"column:token_amount"`     // Number of Tokens withdrawn from Vault (BigInt)
	SharesBurnt  types.BigInt `gorm:"column:shares_burnt"`     // Number of Vault Shares burnt (BigInt)
	Token        *Token       `gorm:"foreignKey:TokenID"`      // Token this Vault will accrue
	TokenID      string       `gorm:"column:token_id"`         // Token ID
	ShareToken   *Token       `gorm:"foreignKey:ShareTokenID"` // Token representing Shares in the Vault
	ShareTokenID string       `gorm:"column:share_token_id"`   // Share Token ID
	SharePrice   types.BigInt `gorm:"column:share_price"`      // Share price (BigInt)
}

func (Withdrawal) TableName() string {
	return "withdrawals"
}

func (w *Withdrawal) Init() {
	w.Timestamp.Zero()
	w.BlockNumber.Zero()

	w.Account = nil
	w.AccountID = ""

	w.Vault = nil
	w.VaultID = ""

	w.TokenAmount.Zero()
	w.SharesBurnt.Zero()

	w.Token = nil
	w.TokenID = ""

	w.ShareToken = nil
	w.ShareTokenID = ""

	w.SharePrice.Zero()
}

func (w *Withdrawal) GetID() string {
	return w.ID
}

func (w *Withdrawal) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, w)
}

func (w *Withdrawal) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, w)
}
