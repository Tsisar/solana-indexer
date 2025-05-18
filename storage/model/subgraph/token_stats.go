package subgraph

import "github.com/Tsisar/solana-indexer/subgraph/types"

type TokenStats struct {
	ID         string           `gorm:"primaryKey;column:id"`         // Aggregated ID (Int8)
	Timestamp  string           `gorm:"column:timestamp;default:0"`   // Timestamp of aggregation
	Vault      *Vault           `gorm:"foreignKey:VaultID"`           // Vault reference
	VaultID    string           `gorm:"column:vault_id"`              // Vault ID
	SharePrice types.BigDecimal `gorm:"column:share_price;default:0"` // Aggregated share price (last) (BigDecimal)
}

func (TokenStats) TableName() string {
	return "token_stats"
}
