package subgraph

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StrategyReportEvent struct {
	ID              string    `gorm:"primaryKey;column:id"`          // The Strategy Report ID
	Timestamp       string    `gorm:"column:timestamp;default:0"`    // Timestamp the strategy report was most recently updated (BigInt)
	BlockNumber     string    `gorm:"column:block_number;default:0"` // Blocknumber the strategy report was most recently updated (BigInt)
	TransactionHash string    `gorm:"column:transaction_hash"`       // Transaction Hash
	Strategy        *Strategy `gorm:"foreignKey:StrategyID"`         // The Strategy reference
	StrategyID      string    `gorm:"column:strategy_id"`            // Strategy ID
	SharePrice      string    `gorm:"column:share_price;default:0"`  // Share price (BigInt)
}

func (*StrategyReportEvent) TableName() string {
	return "strategy_report_events"
}

func (s *StrategyReportEvent) Load(ctx context.Context, db *gorm.DB) (bool, error) {
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

func (s *StrategyReportEvent) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(s).Error
}
