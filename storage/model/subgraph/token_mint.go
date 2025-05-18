package subgraph

import "github.com/Tsisar/solana-indexer/subgraph/types"

type TokenMint struct {
	ID     string           `gorm:"primaryKey;column:id"`    // Mint ID
	To     *Token           `gorm:"foreignKey:ToID"`         // Mint account (ShareToken)
	ToID   string           `gorm:"column:to_id"`            // ShareToken ID
	Amount types.BigDecimal `gorm:"column:amount;default:0"` // Number of Tokens minted (BigDecimal)
}

func (TokenMint) TableName() string {
	return "token_mints"
}
