package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"gorm.io/gorm"
)

type Account struct {
	ID string `gorm:"primaryKey;column:id"` // Account address

	// Derived relationships
	Deposits       []*Deposit              `gorm:"foreignKey:AccountID"`   // From Deposit.account
	Withdrawals    []*Withdrawal           `gorm:"foreignKey:AccountID"`   // From Withdrawal.account
	TokenAccounts  []*TokenWallet          `gorm:"foreignKey:AuthorityID"` // From TokenWallet.authority
	ShareAccounts  []*TokenWallet          `gorm:"foreignKey:AuthorityID"` // From TokenWallet.authority
	VaultPositions []*AccountVaultPosition `gorm:"foreignKey:AccountID"`   // From AccountVaultPosition.account
}

func (Account) TableName() string {
	return "accounts"
}

func (a *Account) Init() {
}

func (a *Account) GetID() string {
	return a.ID
}

func (a *Account) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, a)
}

func (a *Account) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, a)
}
