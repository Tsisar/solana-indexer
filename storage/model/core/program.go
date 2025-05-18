package core

import (
	"time"
)

type Program struct {
	ID        string        `gorm:"primaryKey;column:id"`
	Txns      []Transaction `gorm:"many2many:core.program_transactions;joinForeignKey:program_id;joinReferences:transaction_signature;constraint:OnDelete:CASCADE;" gorm:"column:txns"`
	CreatedAt time.Time     `gorm:"column:created_at;autoCreateTime"`
}

func (Program) TableName() string {
	return "core.programs"
}
