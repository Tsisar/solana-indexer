package subgraph

type TokenBurn struct {
	ID     string `gorm:"primaryKey;column:id"`    // Burn ID
	From   *Token `gorm:"foreignKey:FromID"`       // Burn account (ShareToken)
	FromID string `gorm:"column:from_id"`          // ShareToken ID
	Amount string `gorm:"column:amount;default:0"` // Number of Tokens burnt (BigDecimal)
}

func (TokenBurn) TableName() string {
	return "token_burns"
}
