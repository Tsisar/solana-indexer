package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/monitoring"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type Deposit struct {
	ID           string       `gorm:"primaryKey;column:id"`    // Transaction-Log
	Timestamp    types.BigInt `gorm:"column:timestamp"`        // Timestamp of update (BigInt)
	BlockNumber  types.BigInt `gorm:"column:block_number"`     // Block number of update (BigInt)
	Account      *Account     `gorm:"foreignKey:AccountID"`    // Account making Deposit
	AccountID    string       `gorm:"column:account_id"`       // Account ID
	Vault        *Vault       `gorm:"foreignKey:VaultID"`      // Vault deposited into
	VaultID      string       `gorm:"column:vault_id"`         // Vault ID
	TokenAmount  types.BigInt `gorm:"column:token_amount"`     // Number of Tokens deposited into Vault (BigInt)
	SharesMinted types.BigInt `gorm:"column:shares_minted"`    // Number of new Vault Shares minted (BigInt)
	Token        *Token       `gorm:"foreignKey:TokenID"`      // Token this Vault will accrue
	TokenID      string       `gorm:"column:token_id"`         // Token ID
	ShareToken   *Token       `gorm:"foreignKey:ShareTokenID"` // Token representing Shares in the Vault
	ShareTokenID string       `gorm:"column:share_token_id"`   // Share Token ID
	SharePrice   types.BigInt `gorm:"column:share_price"`      // Share price (BigInt)
}

func (Deposit) TableName() string {
	return "deposits"
}

func (d *Deposit) Init() {
	d.Timestamp.Zero()
	d.BlockNumber.Zero()
	d.Account = nil
	d.AccountID = ""
	d.Vault = nil
	d.VaultID = ""
	d.TokenAmount.Zero()
	d.SharesMinted.Zero()
	d.Token = nil
	d.TokenID = ""
	d.ShareToken = nil
	d.ShareTokenID = ""
	d.SharePrice.Zero()
}

func (d *Deposit) GetID() string {
	return d.ID
}

func (d *Deposit) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, d)
}

func (d *Deposit) Save(ctx context.Context, db *gorm.DB) error {
	monitoring.DepositsTotal.WithLabelValues(d.VaultID, d.TokenID).Inc()
	return generic.Save(ctx, db, d)
}
