package subgraph

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StrategyReport struct {
	ID              string                  `gorm:"primaryKey;column:id"`           // The Strategy Report ID
	Timestamp       string                  `gorm:"column:timestamp;default:0"`     // Timestamp the strategy report was most recently updated (BigInt)
	BlockNumber     string                  `gorm:"column:block_number;default:0"`  // Blocknumber the strategy report was most recently updated (BigInt)
	TransactionHash string                  `gorm:"column:transaction_hash"`        // Transaction Hash
	Strategy        *Strategy               `gorm:"foreignKey:StrategyID"`          // The Strategy reference
	StrategyID      string                  `gorm:"column:strategy_id"`             // Strategy ID
	Gain            string                  `gorm:"column:gain;default:0"`          // Reported gain amount (BigInt)
	Loss            string                  `gorm:"column:loss;default:0"`          // Reported loss amount (BigInt)
	CurrentDebt     string                  `gorm:"column:current_debt;default:0"`  // Reported current debt (BigInt)
	ProtocolFees    string                  `gorm:"column:protocol_fees;default:0"` // Reported protocol fees amount (BigInt)
	TotalFees       string                  `gorm:"column:total_fees;default:0"`    // Reported total fees amount (BigInt)
	TotalShares     string                  `gorm:"column:total_shares;default:0"`  // Reported total shares (BigInt)
	VaultKey        string                  `gorm:"column:vault_key"`               // Vault key
	Results         []*StrategyReportResult `gorm:"foreignKey:CurrentReportID"`     // Results created by this report
}

func (*StrategyReport) TableName() string {
	return "strategy_reports"
}

func (s *StrategyReport) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", s.ID).
		First(s).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (s *StrategyReport) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(s).Error
}
