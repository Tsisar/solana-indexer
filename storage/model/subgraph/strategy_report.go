package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type StrategyReport struct {
	ID              string                  `gorm:"primaryKey;column:id"`       // The Strategy Report ID
	Timestamp       types.BigInt            `gorm:"column:timestamp"`           // Timestamp the strategy report was most recently updated (BigInt)
	BlockNumber     types.BigInt            `gorm:"column:block_number"`        // Blocknumber the strategy report was most recently updated (BigInt)
	TransactionHash string                  `gorm:"column:transaction_hash"`    // Transaction Hash
	Strategy        *Strategy               `gorm:"foreignKey:StrategyID"`      // The Strategy reference
	StrategyID      string                  `gorm:"column:strategy_id"`         // Strategy ID
	Gain            types.BigInt            `gorm:"column:gain"`                // Reported gain amount (BigInt)
	Loss            types.BigInt            `gorm:"column:loss"`                // Reported loss amount (BigInt)
	CurrentDebt     types.BigInt            `gorm:"column:current_debt"`        // Reported current debt (BigInt)
	ProtocolFees    types.BigInt            `gorm:"column:protocol_fees"`       // Reported protocol fees amount (BigInt)
	TotalFees       types.BigInt            `gorm:"column:total_fees"`          // Reported total fees amount (BigInt)
	TotalShares     types.BigInt            `gorm:"column:total_shares"`        // Reported total shares (BigInt)
	VaultKey        string                  `gorm:"column:vault_key"`           // Vault key
	Results         []*StrategyReportResult `gorm:"foreignKey:CurrentReportID"` // Results created by this report
}

func (StrategyReport) TableName() string {
	return "strategy_reports"
}

func (s *StrategyReport) Init() {
	s.Timestamp.Zero()
	s.BlockNumber.Zero()
	s.TransactionHash = ""
	s.Strategy = nil
	s.StrategyID = ""
	s.Gain.Zero()
	s.Loss.Zero()
	s.CurrentDebt.Zero()
	s.ProtocolFees.Zero()
	s.TotalFees.Zero()
	s.TotalShares.Zero()
	s.VaultKey = ""
	s.Results = nil
}

func (s *StrategyReport) GetID() string {
	return s.ID
}

func (s *StrategyReport) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, s)
}

func (s *StrategyReport) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, s)
}
