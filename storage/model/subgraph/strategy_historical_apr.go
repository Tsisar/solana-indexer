package subgraph

type StrategyHistoricalApr struct {
	ID        string `gorm:"primaryKey;column:id"` // Unique ID
	Timestamp string `gorm:"column:timestamp"`     // BigInt
	Apr       string `gorm:"column:apr"`           // BigDecimal

	StrategyID string    `gorm:"column:strategy_id"`    // Foreign key
	Strategy   *Strategy `gorm:"foreignKey:StrategyID"` // Relation to Strategy
}

func (StrategyHistoricalApr) TableName() string {
	return "strategy_historical_aprs"
}
