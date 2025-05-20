package storage

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/core/config"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"gorm.io/gorm/clause"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Gorm wraps a GORM database instance.
type Gorm struct {
	DB *gorm.DB
}

// InitGorm establishes a connection to the PostgreSQL database using configuration values.
// It also ensures the required schemas (e.g., "core") are created.
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

	// Create the core schema if it doesn't exist
	for _, schema := range []string{"core"} {
		if err := db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)).Error; err != nil {
			return nil, fmt.Errorf("failed to create schema %s: %w", schema, err)
		}
	}

	return &Gorm{DB: db}, nil
}

// InitCoreModels runs migrations for the core models used in the indexer,
// sets initial health status, and stores configured program addresses in the database.
func InitCoreModels(ctx context.Context, db *Gorm) error {
	if err := db.DB.AutoMigrate(
		&core.Transaction{},
		&core.Program{},
		&core.Event{},
		&core.IndexerHealth{},
	); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	if err := db.SetHealth(ctx, "unknown", "Just started"); err != nil {
		return fmt.Errorf("failed to set initial health status: %v", err)
	}

	for _, program := range config.App.Programs {
		if err := db.SaveProgram(ctx, program); err != nil {
			return fmt.Errorf("failed to save program address %s: %v", program, err)
		}
	}
	return nil
}

// InitSubgraphModels runs migrations for subgraph-specific models.
// It also creates and validates the `latest_report_id` foreign key.
func InitSubgraphModels(ctx context.Context, db *Gorm) error {
	if err := db.DB.AutoMigrate(
		&subgraph.Meta{},
		&subgraph.BlockInfo{},
		&subgraph.Account{},
		&subgraph.AccountVaultPosition{},
		&subgraph.Accountant{},
		&subgraph.Deposit{},
		&subgraph.DeployFunds{},
		&subgraph.DTFReport{},
		&subgraph.FreeFunds{},
		&subgraph.ShareToken{},
		&subgraph.ShareTokenData{},
		&subgraph.ShareTokenTransfer{},
		&subgraph.Strategy{},
		&subgraph.StrategyHistoricalApr{},
		&subgraph.StrategyReport{},
		&subgraph.StrategyReportEvent{},
		&subgraph.StrategyReportResult{},
		&subgraph.Token{},
		&subgraph.TokenBurn{},
		&subgraph.TokenMint{},
		&subgraph.TokenStats{},
		&subgraph.TokenWallet{},
		&subgraph.Vault{},
		&subgraph.VaultHistoricalApr{},
		&subgraph.Withdrawal{},
		&subgraph.WithdrawalRequest{},
	); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	if err := truncateSubgraphTables(db.DB); err != nil {
		return fmt.Errorf("failed to truncate subgraph tables: %w", err)
	}

	if err := migrateWithLatestReport(db.DB); err != nil {
		return fmt.Errorf("migration latest_report_id failed: %w", err)
	}
	return nil
}

func truncateSubgraphTables(db *gorm.DB) error {
	tables := []string{
		"_meta",
		"_block_info",
		"accounts",
		"account_vault_positions",
		"accountants",
		"deposits",
		"deploy_funds",
		"dtf_reports",
		"free_funds",
		"share_tokens",
		"share_token_data",
		"share_token_transfers",
		"strategies",
		"strategy_historical_aprs",
		"strategy_reports",
		"strategy_report_events",
		"strategy_report_results",
		"tokens",
		"token_burns",
		"token_mints",
		"token_stats",
		"token_wallets",
		"vaults",
		"vault_historical_aprs",
		"withdrawals",
		"withdrawal_requests",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE;", table)).Error; err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
	}

	return nil
}

// migrateWithLatestReport ensures that the `latest_report_id` column and foreign key
// are added to the `strategies` table.
func migrateWithLatestReport(db *gorm.DB) error {
	if !db.Migrator().HasColumn(&subgraph.Strategy{}, "latest_report_id") {
		if err := db.Migrator().AddColumn(&subgraph.Strategy{}, "LatestReportID"); err != nil {
			return fmt.Errorf("add column latest_report_id: %w", err)
		}
	}

	if !db.Migrator().HasConstraint(&subgraph.Strategy{}, "fk_strategies_latest_report") {
		const fk = `
			ALTER TABLE strategies
			  ADD CONSTRAINT fk_strategies_latest_report
			  FOREIGN KEY (latest_report_id)
				REFERENCES strategy_reports(id)
				ON UPDATE CASCADE
				ON DELETE SET NULL;
`
		if err := db.Exec(fk).Error; err != nil {
			return fmt.Errorf("create fk strategies.latest_report_id: %w", err)
		}
	}
	return nil
}

// Close closes the underlying SQL database connection.
func (g *Gorm) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// SetHealth updates the indexer health status (upsert by primary key id = 1).
func (g *Gorm) SetHealth(ctx context.Context, status, reason string) error {
	return g.DB.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "reason"}),
	}).Create(&core.IndexerHealth{
		ID:     1,
		Status: status,
		Reason: reason,
	}).Error
}

// GetHealth retrieves the current indexer health status and reason.
func (g *Gorm) GetHealth(ctx context.Context) (string, string, error) {
	var health core.IndexerHealth
	if err := g.DB.WithContext(ctx).First(&health, 1).Error; err != nil {
		return "", "", err
	}
	return health.Status, health.Reason, nil
}
