package subgraph

type TokenWallet struct {
	ID          string   `gorm:"primaryKey;column:id"`   // Account address (Associated Token Account)
	Authority   *Account `gorm:"foreignKey:AuthorityID"` // Authority
	AuthorityID string   `gorm:"column:authority_id"`    // Authority ID
}

func (TokenWallet) TableName() string {
	return "token_wallets"
}
