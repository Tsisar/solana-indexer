package report

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/internal/subgraph/events"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"github.com/Tsisar/solana-indexer/internal/utils"
	"gorm.io/gorm"
	"math/big"
)

func CreateReport(ctx context.Context, db *gorm.DB, ev events.StrategyReportedEvent, transaction events.Transaction) error {
	log.Infof("[report] Creating report...")
	strategy := subgraph.Strategy{ID: ev.StrategyKey.String()}
	ok, err := strategy.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[report] failed to load strategy: %w", err)
	}

	if ev.Gain.Sign() > 0 || ev.Loss.Sign() > 0 {
		log.Debugf("[report] Generating new report ID for strategy %s.", strategy.ID)

		id := utils.GenerateId(transaction.Signature, ev.StrategyKey.String())
		currentReport := subgraph.StrategyReport{ID: id}
		if _, err := currentReport.Load(ctx, db); err != nil {
			return fmt.Errorf("[report] failed to load strategy report: %w", err)
		}

		currentReport.StrategyID = ev.StrategyKey.String()
		currentReport.BlockNumber = transaction.Slot
		currentReport.Timestamp = transaction.Timestamp
		currentReport.TransactionHash = transaction.Signature
		currentReport.Gain = ev.Gain
		currentReport.Loss = ev.Loss
		currentReport.CurrentDebt = ev.CurrentDebt
		currentReport.ProtocolFees = ev.ProtocolFees
		currentReport.TotalFees = ev.TotalFees
		currentReport.TotalShares = ev.TotalShares
		currentReport.VaultKey = ev.VaultKey.String()

		if err := currentReport.Save(ctx, db); err != nil {
			return fmt.Errorf("[report] failed to save strategy report: %w", err)
		}

		previousReportId := strategy.LatestReportID

		log.Debugf("[report] Getting previous report ID for strategy %s: %v.", strategy.ID, previousReportId)

		strategy.LatestReportID = &currentReport.ID
		if err := strategy.Save(ctx, db); err != nil {
			return fmt.Errorf("[report] failed to save strategy: %w", err)
		}

		if previousReportId != nil && *previousReportId != "" {
			previousReport := subgraph.StrategyReport{ID: *previousReportId}
			ok, err = previousReport.Load(ctx, db)
			if err != nil {
				return fmt.Errorf("[report] failed to load current report: %w", err)
			}
			if !ok {
				log.Warnf("[report] Report result NOT created. Current report not found: %v", previousReportId)
				return nil
			}
			log.Debugf("[report] Creating report result for strategy %s: %v vs %s.", strategy.ID, previousReportId, currentReport.ID)
			if err := createReportResult(ctx, db, previousReport, currentReport, transaction); err != nil {
				return fmt.Errorf("[report] failed to create report result: %w", err)
			}
		} else {
			log.Warnf("[report] Report result NOT created. Previous report not found: %s", previousReportId)
		}
	}
	return nil
}

func createReportResult(ctx context.Context, db *gorm.DB, previousReport subgraph.StrategyReport, currentReport subgraph.StrategyReport, transaction events.Transaction) error {
	log.Infof("[report] Creating report result (latest vs current report)...")
	if currentReport.ID == previousReport.ID {
		log.Warnf("[report] Report result NOT created. Current report is the same as latest report")
		return nil
	}

	id := utils.GenerateId(transaction.Signature, previousReport.ID, currentReport.ID)
	strategyReportResult := subgraph.StrategyReportResult{ID: id}
	if _, err := strategyReportResult.Load(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to load strategy report result: %w", err)
	}

	strategyReportResult.TransactionHash = transaction.Signature
	strategyReportResult.Timestamp = transaction.Timestamp
	strategyReportResult.BlockNumber = transaction.Slot
	strategyReportResult.PreviousReportID = previousReport.ID
	strategyReportResult.CurrentReportID = currentReport.ID
	strategyReportResult.StartTimestamp = previousReport.Timestamp
	strategyReportResult.EndTimestamp = currentReport.Timestamp

	duration := currentReport.Timestamp.Sub(&previousReport.Timestamp)
	msInDays := utils.MillisToDays(duration)

	strategyReportResult.Duration = *duration.ToBigDecimal()

	//TODO: Check if this is correct
	profit := currentReport.Gain.Plus(&currentReport.Loss)

	log.Infof(
		"[report] Report Result - Start / End: %d / %d - Duration: %s (ms) - Profit: %s",
		strategyReportResult.StartTimestamp,
		strategyReportResult.EndTimestamp,
		utils.FormatBigDecimal(&strategyReportResult.Duration, 6),
		utils.FormatBigInt(profit),
	)

	if previousReport.CurrentDebt.Sign() != 0 && msInDays.Sign() != 0 {
		profitOverTotalDebt := profit.ToBigDecimal().SafeDiv(previousReport.CurrentDebt.ToBigDecimal())
		strategyReportResult.DurationPr = *profitOverTotalDebt

		yearOverDuration := utils.DaysToYearFactor(msInDays)
		hundred := &types.BigDecimal{Float: big.NewFloat(100)}

		apr := profitOverTotalDebt.Mul(yearOverDuration).Mul(hundred)
		strategyReportResult.Apr = *apr

		log.Infof(
			"[report] Report Result - Duration: %s ms / %s days - Duration (Year): %s - Profit / Total Debt: %s / APR: %s - TxHash: %s",
			utils.FormatBigInt(duration),                   // milliseconds
			utils.FormatBigDecimal(msInDays, 6),            // days
			utils.FormatBigDecimal(yearOverDuration, 6),    // year factor
			utils.FormatBigDecimal(profitOverTotalDebt, 6), // ratio
			utils.FormatBigDecimal(apr, 6),                 // %
			transaction.Signature,
		)
	}

	strategy := subgraph.Strategy{ID: previousReport.StrategyID}
	ok, err := strategy.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[report] failed to load strategy: %w", err)
	}

	vault := subgraph.Vault{ID: strategy.VaultID}
	ok, err = vault.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[report] failed to load vault: %w", err)
	}

	reportCount := strategy.ReportsCount
	numeratorVault := vault.Apr.Plus(&strategyReportResult.Apr)
	numeratorStrategy := strategy.Apr.Plus(&strategyReportResult.Apr)

	var vaultApr *types.BigDecimal
	var strategyApr *types.BigDecimal

	if reportCount.Sign() == 0 {
		vaultApr = numeratorVault
		strategyApr = numeratorStrategy
	} else {
		vaultApr = numeratorVault.SafeDiv(reportCount.ToBigDecimal())
		strategyApr = numeratorStrategy.SafeDiv(reportCount.ToBigDecimal())
	}

	vault.Apr = *vaultApr
	strategy.Apr = *strategyApr

	newVaultHistoricalApr := subgraph.VaultHistoricalApr{ID: id}
	if _, err := newVaultHistoricalApr.Load(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to load vault historical APR: %w", err)
	}
	newVaultHistoricalApr.Timestamp = transaction.Timestamp
	newVaultHistoricalApr.Apr = *vaultApr
	newVaultHistoricalApr.VaultID = vault.ID
	if err := newVaultHistoricalApr.Save(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to save vault historical APR: %w", err)
	}

	newStrategyHistoricalApr := subgraph.StrategyHistoricalApr{ID: id}
	if _, err := newStrategyHistoricalApr.Load(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to load strategy historical APR: %w", err)
	}
	newStrategyHistoricalApr.Timestamp = transaction.Timestamp
	newStrategyHistoricalApr.Apr = *strategyApr
	newStrategyHistoricalApr.StrategyID = strategy.ID
	if err := newStrategyHistoricalApr.Save(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to save strategy historical APR: %w", err)
	}

	strategy.ReportsCount = *reportCount.Plus(&types.BigInt{Int: big.NewInt(1)})
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to save strategy: %w", err)
	}
	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to save vault: %w", err)
	}
	if err := strategyReportResult.Save(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to save strategy report result: %w", err)
	}

	return nil
}

func CreateReportEvent(ctx context.Context, db *gorm.DB, ev events.StrategyReportedEvent, transaction events.Transaction) error {
	log.Infof("[report] Creating report event...")

	id := utils.GenerateId(transaction.Signature, ev.StrategyKey.String())
	strategyReportEvent := subgraph.StrategyReportEvent{ID: id}
	if _, err := strategyReportEvent.Load(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to load strategy report event: %w", err)
	}

	strategyReportEvent.TransactionHash = transaction.Signature
	strategyReportEvent.StrategyID = ev.StrategyKey.String()
	strategyReportEvent.BlockNumber = transaction.Slot
	strategyReportEvent.Timestamp = transaction.Timestamp
	strategyReportEvent.SharePrice = ev.SharePrice

	if err := strategyReportEvent.Save(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to save strategy report event: %w", err)
	}
	return nil
}

func CreateShareTokenData(ctx context.Context, db *gorm.DB, ev events.StrategyReportedEvent, transaction events.Transaction) error {
	log.Infof("[report] Creating share token data...")
	strategy := subgraph.Strategy{ID: ev.StrategyKey.String()}
	ok, err := strategy.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[report] failed to load strategy: %w", err)
	}

	vault := subgraph.Vault{ID: strategy.VaultID}
	ok, err = vault.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[report] failed to load vault: %w", err)
	}

	id := utils.GenerateId(transaction.Signature, strategy.VaultID, transaction.Timestamp.String())
	shareTokenData := subgraph.ShareTokenData{
		ID:         id,
		VaultID:    strategy.VaultID,
		Timestamp:  transaction.Timestamp,
		SharePrice: ev.SharePrice,
	}
	if err := shareTokenData.Save(ctx, db); err != nil {
		return fmt.Errorf("[report] failed to save share token data: %w", err)
	}

	return nil
}
