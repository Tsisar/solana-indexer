package subgraph

type TokenMint struct {
	ID     string `gorm:"primaryKey;column:id"`    // Mint ID
	To     *Token `gorm:"foreignKey:ToID"`         // Mint account (ShareToken)
	ToID   string `gorm:"column:to_id"`            // ShareToken ID
	Amount string `gorm:"column:amount;default:0"` // Number of Tokens minted (BigDecimal)
}

func (TokenMint) TableName() string {
	return "token_mints"
}
