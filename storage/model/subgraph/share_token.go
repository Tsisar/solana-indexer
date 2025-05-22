package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type ShareToken struct {
	ID          string                `gorm:"primaryKey;column:id"` // Associated Token Account
	Mint        []*TokenMint          `gorm:"foreignKey:ToID"`      // Token mint account
	Burn        []*TokenBurn          `gorm:"foreignKey:FromID"`    // Token burn account
	TransferIn  []*ShareTokenTransfer `gorm:"foreignKey:ToID"`      // Token transfer account in
	TransferOut []*ShareTokenTransfer `gorm:"foreignKey:FromID"`    // Token transfer account out

	TotalMinted      types.BigDecimal `gorm:"column:total_minted"`       // Total Minted (BigDecimal)
	TotalBurnt       types.BigDecimal `gorm:"column:total_burnt"`        // Total Burnt (BigDecimal)
	TotalTransferIn  types.BigDecimal `gorm:"column:total_transfer_in"`  // Total Transfer In (BigDecimal)
	TotalTransferOut types.BigDecimal `gorm:"column:total_transfer_out"` // Total Transfer Out (BigDecimal)
	CurrentPrice     types.BigInt     `gorm:"column:current_price"`      // Current price of the Token (BigInt)
}

func (ShareToken) TableName() string {
	return "share_tokens"
}

func (s *ShareToken) Init() {
	s.TotalMinted.Zero()
	s.TotalBurnt.Zero()
	s.TotalTransferIn.Zero()
	s.TotalTransferOut.Zero()
	s.CurrentPrice.Zero()
}

func (s *ShareToken) GetID() string {
	return s.ID
}

func (s *ShareToken) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, s)
}

func (s *ShareToken) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, s)
}
