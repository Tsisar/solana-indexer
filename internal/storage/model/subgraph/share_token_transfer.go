package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"gorm.io/gorm"
)

type ShareTokenTransfer struct {
	ID          string           `gorm:"primaryKey;column:id"` // ID
	Mint        *Token           `gorm:"foreignKey:MintID"`    // Mint account (Token)
	MintID      string           `gorm:"column:mint_id"`       // Token mint address
	AuthorityID string           `gorm:"column:authority_id"`  // Authority ID (Int8)
	To          *ShareToken      `gorm:"foreignKey:ToID"`      // Transfer to account (ShareToken)
	ToID        string           `gorm:"column:to_id"`         // ShareToken ID (recipient)
	From        *ShareToken      `gorm:"foreignKey:FromID"`    // Transfer from account (ShareToken)
	FromID      string           `gorm:"column:from_id"`       // ShareToken ID (sender)
	Amount      types.BigDecimal `gorm:"column:amount"`        // Number of Tokens transferred (BigDecimal)
}

func (ShareTokenTransfer) TableName() string {
	return "share_token_transfers"
}

func (s *ShareTokenTransfer) Init() {
	s.AuthorityID = ""
	s.To = nil
	s.ToID = ""
	s.From = nil
	s.FromID = ""
	s.Amount.Zero()
}

func (s *ShareTokenTransfer) GetID() string {
	return s.ID
}

func (s *ShareTokenTransfer) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, s)
}

func (s *ShareTokenTransfer) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, s)
}
