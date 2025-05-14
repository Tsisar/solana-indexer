package storage

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/internal/config"
	"gorm.io/gorm/clause"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Gorm struct {
	DB *gorm.DB
}

type Program struct {
	ID        string        `gorm:"primaryKey"`
	Txns      []Transaction `gorm:"many2many:program_transactions;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time     `gorm:"autoCreateTime"`
}

func InitGorm() (*Gorm, error) {
	cfg := config.App.Postgres
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=UTC",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get raw db failed: %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	if err := db.AutoMigrate(&Transaction{}, &Program{}, &Event{}); err != nil {
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return &Gorm{DB: db}, nil
}

func (g *Gorm) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

func (g *Gorm) SaveProgram(ctx context.Context, address string) error {
	return g.DB.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&Program{
			ID: address,
		}).Error
}
