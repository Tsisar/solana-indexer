package subgraph

type DTFReport struct {
	ID          string `gorm:"primaryKey;column:id"`          // The Strategy Report ID
	TotalAssets string `gorm:"column:total_assets;default:0"` // Total amount of assets deposited in strategies (BigInt)
	Timestamp   string `gorm:"column:timestamp;default:0"`    // Timestamp the strategy report was most recently updated (BigInt)
}

func (DTFReport) TableName() string {
	return "dtf_reports"
}
