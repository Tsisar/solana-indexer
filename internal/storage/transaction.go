package storage

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/internal/config"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Transaction struct {
	Signature string `gorm:"primaryKey"`
	Slot      uint64
	BlockTime int64
	JsonTx    datatypes.JSON `gorm:"type:jsonb"`
	Parsed    bool           `gorm:"default:false"`
	Programs  []Program      `gorm:"many2many:program_transactions;constraint:OnDelete:CASCADE;"`
	Events    []Event        `gorm:"foreignKey:TransactionSignature;references:Signature;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
}

func (g *Gorm) SaveTransaction(ctx context.Context, signature string, slot uint64) error {
	transactions := Transaction{
		Signature: signature,
		Slot:      slot,
	}

	return g.DB.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&transactions).Error
}

func (g *Gorm) UpdateTransactionRaw(ctx context.Context, signature string, raw []byte) error {
	return g.DB.WithContext(ctx).
		Model(&Transaction{}).
		Where("signature = ?", signature).
		Updates(map[string]interface{}{
			"json_tx": datatypes.JSON(raw),
		}).Error
}

func (g *Gorm) MarkParsed(ctx context.Context, signature string) error {
	return g.DB.WithContext(ctx).
		Model(&Transaction{}).
		Where("signature = ?", signature).
		Update("parsed", true).Error
}

func (g *Gorm) IsParsed(ctx context.Context, signature string) (bool, error) {
	var count int64
	if err := g.DB.WithContext(ctx).
		Model(&Transaction{}).
		Where("signature = ? AND parsed = true", signature).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if transaction is parsed: %w", err)
	}

	return count > 0, nil
}

func (g *Gorm) AssociateTransactionWithProgram(ctx context.Context, signature, programID string) error {
	var tx Transaction
	if err := g.DB.WithContext(ctx).First(&tx, "signature = ?", signature).Error; err != nil {
		return fmt.Errorf("transaction not found: %w", err)
	}

	var prog Program
	if err := g.DB.WithContext(ctx).First(&prog, "id = ?", programID).Error; err != nil {
		return fmt.Errorf("program not found: %w", err)
	}

	if err := g.DB.WithContext(ctx).Model(&tx).Association("Programs").Append(&prog); err != nil {
		return fmt.Errorf("failed to associate: %w", err)
	}

	return nil
}

func (g *Gorm) IsReady(ctx context.Context) (bool, error) {
	var count int64
	if err := g.DB.WithContext(ctx).
		Model(&Transaction{}).
		Where("parsed = false").
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to count unparsed transactions: %w", err)
	}

	return count == 0, nil
}

func (g *Gorm) GetOrderedNoParsedSignatures(ctx context.Context) ([]string, error) {
	addresses := config.App.Programs
	var signatures []string

	err := g.DB.WithContext(ctx).
		Raw(`
		SELECT transactions.signature
		FROM transactions
		JOIN program_transactions ON program_transactions.transaction_signature = transactions.signature
		JOIN programs ON programs.id = program_transactions.program_id
		WHERE programs.id IN ?
		  AND transactions.parsed = false
		GROUP BY transactions.signature, transactions.block_time
		ORDER BY transactions.block_time ASC
	`, addresses).
		Scan(&signatures).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch unparsed signatures: %w", err)
	}

	return signatures, nil
}

func (g *Gorm) GetOrderedNoRawSignatures(ctx context.Context) ([]string, error) {
	addresses := config.App.Programs
	var signatures []string

	err := g.DB.WithContext(ctx).
		Raw(`
		SELECT transactions.signature
		FROM transactions
		JOIN program_transactions ON program_transactions.transaction_signature = transactions.signature
		JOIN programs ON programs.id = program_transactions.program_id
		WHERE programs.id IN ?
		  AND transactions.json_tx IS NULL
		GROUP BY transactions.signature, transactions.block_time
		ORDER BY transactions.block_time ASC
	`, addresses).
		Scan(&signatures).Error

	if err != nil {
		return nil, fmt.Errorf("failed to fetch unparsed signatures: %w", err)
	}

	return signatures, nil
}

func (g *Gorm) GetLatestSavedSignature(ctx context.Context, programID string) (string, error) {
	var signature string
	query := `
		SELECT t.signature
		FROM transactions t
		JOIN program_transactions pt ON pt.transaction_signature = t.signature
		WHERE pt.program_id = ?
		ORDER BY t.block_time DESC
		LIMIT 1`

	err := g.DB.WithContext(ctx).
		Raw(query, programID).
		Scan(&signature).Error
	if err != nil {
		return "", fmt.Errorf("failed to fetch last saved signature: %w", err)
	}
	return signature, nil
}

func (g *Gorm) GetRawTransaction(ctx context.Context, signature string) ([]byte, error) {
	var tx Transaction
	err := g.DB.WithContext(ctx).
		First(&tx, "signature = ?", signature).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return tx.JsonTx, nil
}
