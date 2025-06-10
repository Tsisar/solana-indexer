package core

import (
	"context"
	"github.com/Tsisar/solana-indexer/internal/config"
	"github.com/Tsisar/solana-indexer/internal/storage/model/generic"
	"gorm.io/gorm"
	"time"
)

type Status struct {
	ID           string    `gorm:"primaryKey;column:id"`
	CurrentBlock uint64    `gorm:"column:current_block"`
	Synced       bool      `gorm:"column:synced"`
	HasError     bool      `gorm:"column:has_error"`
	ErrorMsg     string    `gorm:"column:error_msg"`
	Version      string    `gorm:"column:version"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

func (Status) TableName() string {
	return "core.status"
}

func (s *Status) Init() {
	s.ID = "1" // Default ID for the status record
	s.CurrentBlock = 0
	s.Synced = false
	s.HasError = false
	s.ErrorMsg = ""
	s.Version = config.App.Version
	s.UpdatedAt = time.Now()
	s.CreatedAt = time.Now()
}

func (s *Status) GetID() string {
	return s.ID
}

func (s *Status) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	s.ID = "1" // Ensure the ID is set to the default value
	return generic.Load(ctx, db, s)
}

func (s *Status) Save(ctx context.Context, db *gorm.DB) error {
	s.ID = "1" // Ensure the ID is set to the default value
	s.HasError = s.ErrorMsg != ""
	s.UpdatedAt = time.Now()
	return generic.Save(ctx, db, s)
}
