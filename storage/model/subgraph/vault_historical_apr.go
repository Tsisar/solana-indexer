package subgraph

type VaultHistoricalApr struct {
	ID        string `gorm:"primaryKey;column:id"`       // ID
	Timestamp string `gorm:"column:timestamp;default:0"` // Time in UTC (BigInt)
	Apr       string `gorm:"column:apr;default:0"`       // The Annual Percentage Rate (BigDecimal)
	Vault     *Vault `gorm:"foreignKey:VaultID"`         // The Vault
	VaultID   string `gorm:"column:vault_id"`            // Vault ID
}

func (VaultHistoricalApr) TableName() string {
	return "vault_historical_aprs"
}
