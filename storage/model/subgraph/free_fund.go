package subgraph

type FreeFunds struct {
	ID         string    `gorm:"primaryKey;column:id"`       // The Strategy Fund ID
	Strategy   *Strategy `gorm:"foreignKey:StrategyID"`      // The Strategy
	StrategyID string    `gorm:"column:strategy_id"`         // Strategy ID
	Amount     string    `gorm:"column:amount;default:0"`    // Total amount of assets deposited in strategies (BigInt)
	Timestamp  string    `gorm:"column:timestamp;default:0"` // Timestamp the strategy report was most recently updated (BigInt)
}

func (FreeFunds) TableName() string {
	return "free_funds"
}
