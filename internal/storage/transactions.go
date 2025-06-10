package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/Tsisar/solana-indexer/internal/config"
	"github.com/Tsisar/solana-indexer/internal/storage/model/core"
	"github.com/gagliardetto/solana-go"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SaveTransaction inserts a transaction if it doesn't exist, and associates it with a program via M2M.
// It uses ON CONFLICT DO NOTHING to avoid duplicates.
func (g *Gorm) SaveTransaction(ctx context.Context, tx *core.Transaction, programID string) error {
	// Insert the transaction if not already present
	if err := g.DB.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(tx).Error; err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	// Load the program by ID
	var prog core.Program
	if err := g.DB.WithContext(ctx).
		First(&prog, "id = ?", programID).Error; err != nil {
		return fmt.Errorf("program not found: %w", err)
	}

	// Associate transaction with the program (many-to-many)
	if err := g.DB.WithContext(ctx).
		Model(&core.Transaction{Signature: tx.Signature}).
		Association("Programs").
		Append(&prog); err != nil {
		return fmt.Errorf("failed to associate transaction with program: %w", err)
	}

	return nil
}

// AssociateTransactionWithProgram links a transaction to a program via the many-to-many relationship.
func (g *Gorm) AssociateTransactionWithProgram(ctx context.Context, signature, programID string) error {
	var prog core.Program
	if err := g.DB.WithContext(ctx).
		First(&prog, "id = ?", programID).Error; err != nil {
		return fmt.Errorf("program not found: %w", err)
	}

	if err := g.DB.WithContext(ctx).
		Model(&core.Transaction{Signature: signature}).
		Association("Programs").
		Append(&prog); err != nil {
		return fmt.Errorf("failed to associate transaction with program: %w", err)
	}

	return nil
}

// UpdateTransactionRaw updates the `json_tx` field of a transaction by its signature.
func (g *Gorm) UpdateTransactionRaw(ctx context.Context, signature string, raw []byte) error {
	return g.DB.WithContext(ctx).
		Model(&core.Transaction{}).
		Where("signature = ?", signature).
		Updates(map[string]interface{}{
			"json_tx": datatypes.JSON(raw),
		}).Error
}

// MarkParsed sets the `parsed` flag of a transaction to true.
func (g *Gorm) MarkParsed(ctx context.Context, signature string) error {
	return g.DB.WithContext(ctx).
		Model(&core.Transaction{}).
		Where("signature = ?", signature).
		Update("parsed", true).Error
}

// IsParsed checks whether a transaction has already been parsed.
func (g *Gorm) IsParsed(ctx context.Context, signature string) (bool, error) {
	var count int64
	if err := g.DB.WithContext(ctx).
		Model(&core.Transaction{}).
		Where("signature = ? AND parsed = true", signature).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if transaction is parsed: %w", err)
	}
	return count > 0, nil
}

// IsRawFetched checks whether a transaction has already been fetched in raw format.
func (g *Gorm) IsRawFetched(ctx context.Context, signature string) (bool, error) {
	var count int64
	if err := g.DB.WithContext(ctx).
		Model(&core.Transaction{}).
		Where("signature = ? AND json_tx IS NOT NULL", signature).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if transaction is fetched: %w", err)
	}
	return count > 0, nil
}

// IsReady returns true if there are no unparsed transactions left.
func (g *Gorm) IsReady(ctx context.Context) (bool, error) {
	var count int64
	if err := g.DB.WithContext(ctx).
		Model(&core.Transaction{}).
		Where("parsed = false").
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to count unparsed transactions: %w", err)
	}
	return count == 0, nil
}

// GetOrderedNoParsedSignatures returns signatures of transactions (optionally only unparsed)
// that are associated with configured programs, ordered by block_time.
func (g *Gorm) GetOrderedNoParsedSignatures(ctx context.Context) ([]string, error) {
	addresses := config.App.Programs
	var signatures []string

	query := `
		SELECT t.signature
		FROM core.transactions t
		JOIN core.program_transactions pt ON pt.transaction_signature = t.signature
		JOIN core.programs p ON p.id = pt.program_id
		WHERE p.id IN ?`

	args := []any{addresses}

	query += `
		AND t.parsed = false
		GROUP BY t.signature, t.block_time
		ORDER BY t.block_time ASC`

	err := g.DB.WithContext(ctx).
		Raw(query, args...).
		Scan(&signatures).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch signatures: %w", err)
	}

	return signatures, nil
}

// GetOrderedNoRawSignatures returns signatures of transactions that are missing raw JSON payloads.
func (g *Gorm) GetOrderedNoRawSignatures(ctx context.Context) ([]string, error) {
	addresses := config.App.Programs
	var signatures []string

	err := g.DB.WithContext(ctx).
		Raw(`
		SELECT t.signature
		FROM core.transactions t
		JOIN core.program_transactions pt ON pt.transaction_signature = t.signature
		JOIN core.programs p ON p.id = pt.program_id
		WHERE p.id IN ?
		  AND t.json_tx IS NULL
		GROUP BY t.signature, t.block_time
		ORDER BY t.block_time ASC
	`, addresses).
		Scan(&signatures).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions without raw: %w", err)
	}
	return signatures, nil
}

// GetLatestSavedSignature returns the latest saved signature (by block_time) for a given program ID.
func (g *Gorm) GetLatestSavedSignature(ctx context.Context, programID string) (solana.Signature, error) {
	var signature string
	query := `
		SELECT t.signature
		FROM core.transactions t
		JOIN core.program_transactions pt ON pt.transaction_signature = t.signature
		WHERE pt.program_id = ?
		ORDER BY t.block_time DESC
		LIMIT 1`

	err := g.DB.WithContext(ctx).
		Raw(query, programID).
		Scan(&signature).Error
	if err != nil {
		return solana.Signature{}, fmt.Errorf("failed to fetch last saved signature: %w", err)
	}
	if signature != "" {
		return solana.MustSignatureFromBase58(signature), nil
	}

	return solana.Signature{}, nil
}

// GetRawTransaction returns the raw JSON transaction bytes for a given signature.
// Returns nil if the transaction is not found.
func (g *Gorm) GetRawTransaction(ctx context.Context, signature string) ([]byte, error) {
	var tx core.Transaction
	err := g.DB.WithContext(ctx).
		First(&tx, "signature = ?", signature).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch transaction: %w", err)
	}
	return tx.JsonTx, nil
}

func MarkAllTransactionsUnparsed(ctx context.Context, db *gorm.DB) error {
	result := db.WithContext(ctx).
		Model(&core.Transaction{}).
		Where("parsed = ?", true).
		Update("parsed", false)

	if result.Error != nil {
		return fmt.Errorf("failed to mark transactions as unparsed: %w", result.Error)
	}

	return nil
}
