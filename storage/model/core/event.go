package core

import (
	"gorm.io/datatypes"
	"time"
)

type Event struct {
	TransactionSignature string         `gorm:"column:transaction_signature;primaryKey"`
	LogIndex             int            `gorm:"column:log_index;primaryKey"`
	BlockTime            int64          `gorm:"column:block_time"`
	Slot                 uint64         `gorm:"column:slot"`
	Name                 string         `gorm:"column:name"`
	JsonEv               datatypes.JSON `gorm:"column:json_ev;type:jsonb"`
	Mapped               bool           `gorm:"column:mapped"`
	CreatedAt            time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt            time.Time      `gorm:"column:updated_at;autoUpdateTime"`
}

func (*Event) TableName() string {
	return "core.events"
}
