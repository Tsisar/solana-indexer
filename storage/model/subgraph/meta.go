package subgraph

import (
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type BlockInfo struct {
	ID         uint   `gorm:"primaryKey;column:id"`
	Hash       string `gorm:"column:hash"`
	Number     uint64 `gorm:"column:number"`
	ParentHash string `gorm:"column:parent_hash"`
	Timestamp  int64  `gorm:"column:timestamp"`
}

func (BlockInfo) TableName() string {
	return "_block_info"
}

func (b *BlockInfo) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(b).Error
}

type Meta struct {
	ID                uint       `gorm:"primaryKey;autoIncrement;column:id"`
	Deployment        string     `gorm:"column:deployment"`
	HasIndexingErrors bool       `gorm:"column:has_indexing_errors"`
	ErrorMessage      string     `gorm:"column:error_message"`
	Block             *BlockInfo `gorm:"foreignKey:BlockID"`
	BlockID           uint       `gorm:"column:block_id"`
	CreatedAt         time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (*Meta) TableName() string {
	return "_meta"
}

func (m *Meta) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(m).Error
}
