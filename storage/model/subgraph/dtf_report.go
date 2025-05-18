package subgraph

import "github.com/Tsisar/solana-indexer/subgraph/types"

type DTFReport struct {
	ID          string       `gorm:"primaryKey;column:id"`          // The Strategy Report ID
	TotalAssets types.BigInt `gorm:"column:total_assets;default:0"` // Total amount of assets deposited in strategies (BigInt)
	Timestamp   types.BigInt `gorm:"column:timestamp;default:0"`    // Timestamp the strategy report was most recently updated (BigInt)
}

func (DTFReport) TableName() string {
	return "dtf_reports"
}
