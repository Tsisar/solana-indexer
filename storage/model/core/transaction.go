package core

import (
	"gorm.io/datatypes"
	"time"
)

type Transaction struct {
	Signature string         `gorm:"primaryKey;column:signature"`
	Slot      uint64         `gorm:"column:slot"`
	BlockTime int64          `gorm:"column:block_time"`
	JsonTx    datatypes.JSON `gorm:"column:json_tx;type:jsonb"`
	Parsed    bool           `gorm:"column:parsed;default:false"`
	Programs  []Program      `gorm:"many2many:core.program_transactions;joinForeignKey:transaction_signature;joinReferences:program_id;constraint:OnDelete:CASCADE;" gorm:"column:programs"`
	Events    []Event        `gorm:"foreignKey:TransactionSignature;references:Signature;constraint:OnDelete:CASCADE" gorm:"column:events"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
}

func (Transaction) TableName() string {
	return "core.transactions"
}
