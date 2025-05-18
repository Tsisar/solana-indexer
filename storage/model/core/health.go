package core

import "time"

type IndexerHealth struct {
	ID        uint      `gorm:"primaryKey;autoIncrement;column:id"`
	Status    string    `gorm:"column:status;not null;default:healthy"`
	Reason    string    `gorm:"column:reason;type:text"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (IndexerHealth) TableName() string {
	return "core.indexer_health"
}
