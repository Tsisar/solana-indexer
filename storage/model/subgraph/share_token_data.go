package subgraph

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type ShareTokenData struct {
	ID         string       `gorm:"primaryKey;column:id"` // ID (Int8)
	Vault      *Vault       `gorm:"foreignKey:VaultID"`   // Reference to Vault
	VaultID    string       `gorm:"column:vault_id"`      // Vault ID
	Timestamp  types.BigInt `gorm:"column:timestamp"`     // Timestamp of the record
	SharePrice types.BigInt `gorm:"column:share_price"`   // Share price (BigInt)
}

func (ShareTokenData) TableName() string {
	return "share_token_data"
}

func (s *ShareTokenData) Init() {
	s.Vault = nil
	s.VaultID = ""
	s.Timestamp.Zero()
	s.SharePrice.Zero()
}

func (s *ShareTokenData) GetID() string {
	return s.ID
}

func (s *ShareTokenData) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.Load(ctx, db, s)
}

func (s *ShareTokenData) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, s)
}
