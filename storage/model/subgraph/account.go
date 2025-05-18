package subgraph

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
