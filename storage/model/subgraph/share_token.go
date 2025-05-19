package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShareToken struct {
	ID          string                `gorm:"primaryKey;column:id"` // Associated Token Account
	Mint        []*TokenMint          `gorm:"foreignKey:ToID"`      // Token mint account
	Burn        []*TokenBurn          `gorm:"foreignKey:FromID"`    // Token burn account
	TransferIn  []*ShareTokenTransfer `gorm:"foreignKey:ToID"`      // Token transfer account in
	TransferOut []*ShareTokenTransfer `gorm:"foreignKey:FromID"`    // Token transfer account out

	TotalMinted      types.BigDecimal `gorm:"column:total_minted;default:0"`       // Total Minted (BigDecimal)
	TotalBurnt       types.BigDecimal `gorm:"column:total_burnt;default:0"`        // Total Burnt (BigDecimal)
	TotalTransferIn  types.BigDecimal `gorm:"column:total_transfer_in;default:0"`  // Total Transfer In (BigDecimal)
	TotalTransferOut types.BigDecimal `gorm:"column:total_transfer_out;default:0"` // Total Transfer Out (BigDecimal)
	CurrentPrice     types.BigInt     `gorm:"column:current_price;default:0"`      // Current price of the Token (BigInt)
}

func (ShareToken) TableName() string {
	return "share_tokens"
}

func (s *ShareToken) Load(ctx context.Context, db *gorm.DB) (bool, error) {
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

func (s *ShareToken) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(s).Error
}
