package subgraph

import "github.com/Tsisar/solana-indexer/subgraph/types"

type ShareTokenTransfer struct {
	ID     string           `gorm:"primaryKey;column:id"`    // ID
	To     *Token           `gorm:"foreignKey:ToID"`         // Transfer to account (ShareToken)
	ToID   string           `gorm:"column:to_id"`            // ShareToken ID (recipient)
	From   *Token           `gorm:"foreignKey:FromID"`       // Transfer from account (ShareToken)
	FromID string           `gorm:"column:from_id"`          // ShareToken ID (sender)
	Amount types.BigDecimal `gorm:"column:amount;default:0"` // Number of Tokens transferred (BigDecimal)
}

func (ShareTokenTransfer) TableName() string {
	return "share_token_transfers"
}
