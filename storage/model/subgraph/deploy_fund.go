package subgraph

import "github.com/Tsisar/solana-indexer/subgraph/types"

type DeployFunds struct {
	ID         string       `gorm:"primaryKey;column:id"`       // The Strategy Deploy ID
	Strategy   *Strategy    `gorm:"foreignKey:StrategyID"`      // The Strategy
	StrategyID string       `gorm:"column:strategy_id"`         // Strategy ID
	Amount     types.BigInt `gorm:"column:amount;default:0"`    // Total amount of assets deposited in strategy (BigInt)
	Timestamp  types.BigInt `gorm:"column:timestamp;default:0"` // Timestamp the strategy report was most recently updated (BigInt)
}

func (DeployFunds) TableName() string {
	return "deploy_funds"
}
