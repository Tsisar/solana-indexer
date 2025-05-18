package report

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/events"
	coremodel "github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/Tsisar/solana-indexer/subgraph/model"
	"github.com/Tsisar/solana-indexer/subgraph/utils"
	"gorm.io/gorm"
	"math/big"
	"strconv"
)

func CreateReport(ctx context.Context, db *gorm.DB, event coremodel.Event, ev events.StrategyReportedEvent) error {
	log.Infof("[createReport] Creating report...")
	strategy := modelsss.Strategy{ID: ev.StrategyKey.String()}
	ok, err := strategy.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[createReport] failed to load strategy: %w", err)
	}

	if ev.Gain > 0 || ev.Loss > 0 {
		log.Debugf("[createReport] Generating new report ID for strategy %s.", strategy.ID)

		id := utils.GenerateId(event)
		currentReport := modelsss.StrategyReport{ID: id}
		if _, err := currentReport.Load(ctx, db); err != nil {
			return fmt.Errorf("[createReport] failed to load strategy report: %w", err)
		}

		currentReport.StrategyID = ev.StrategyKey.String()
		currentReport.BlockNumber = strconv.FormatUint(event.Slot, 10)
		currentReport.Timestamp = strconv.FormatInt(ev.Timestamp, 10)
		currentReport.TransactionHash = event.TransactionSignature
		currentReport.Gain = strconv.FormatUint(ev.Gain, 10)
		currentReport.Loss = strconv.FormatUint(ev.Loss, 10)
		currentReport.CurrentDebt = strconv.FormatUint(ev.CurrentDebt, 10)
		currentReport.ProtocolFees = strconv.FormatUint(ev.ProtocolFees, 10)
		currentReport.TotalFees = strconv.FormatUint(ev.TotalFees, 10)
		currentReport.TotalShares = strconv.FormatUint(ev.TotalShares, 10)
		currentReport.VaultKey = ev.VaultKey.String()

		if err := currentReport.Save(ctx, db); err != nil {
			return fmt.Errorf("[createReport] failed to save strategy report: %w", err)
		}

		previousReportID := strategy.LatestReportID
		log.Debugf("[createReport] Getting previous report ID for strategy %s: %s.", strategy.ID, previousReportID)

		strategy.LatestReportID = &currentReport.ID
		if err := strategy.Save(ctx, db); err != nil {
			return fmt.Errorf("[createReport] failed to save strategy: %w", err)
		}

		if previousReportID != nil {
			previousReport := modelsss.StrategyReport{ID: *previousReportID}
			ok, err = previousReport.Load(ctx, db)
			if err != nil {
				return fmt.Errorf("[createReport] failed to load current report: %w", err)
			}
			if !ok {
				log.Warnf("[createReport] Report result NOT created. Current report not found: %s", previousReportID)
				return nil
			}
			if err := createReportResult(ctx, db, event, previousReport, currentReport); err != nil {
				return fmt.Errorf("[createReport] failed to create report result: %w", err)
			}
		}
	}
	return nil
}

func CreateReportEvent(ctx context.Context, db *gorm.DB, event coremodel.Event, ev events.StrategyReportedEvent) error {
	log.Infof("[createReportEvent] Creating report event...")

	id := utils.GenerateId(event)
	strategyReportEvent := modelsss.StrategyReportEvent{ID: id}
	if _, err := strategyReportEvent.Load(ctx, db); err != nil {
		return fmt.Errorf("[createReportEvent] failed to load strategy report event: %w", err)
	}

	strategyReportEvent.TransactionHash = event.TransactionSignature
	strategyReportEvent.StrategyID = ev.StrategyKey.String()
	strategyReportEvent.BlockNumber = strconv.FormatUint(event.Slot, 10)
	strategyReportEvent.Timestamp = strconv.FormatInt(ev.Timestamp, 10)
	strategyReportEvent.SharePrice = strconv.FormatUint(ev.SharePrice, 10)

	if err := strategyReportEvent.Save(ctx, db); err != nil {
		return fmt.Errorf("[createReportEvent] failed to save strategy report event: %w", err)
	}
	return nil
}

func CreateShareTokenData(ctx context.Context, db *gorm.DB, event coremodel.Event, ev events.StrategyReportedEvent) error {
	log.Infof("[createShareTokenData] Creating share token data...")
	strategy := modelsss.Strategy{ID: ev.StrategyKey.String()}
	ok, err := strategy.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[createShareTokenData] failed to load strategy: %w", err)
	}

	shareTokenData := modelsss.ShareTokenData{
		ID:         "0",
		VaultID:    strategy.VaultID,
		Timestamp:  strconv.FormatInt(event.BlockTime, 10),
		SharePrice: strconv.FormatUint(ev.SharePrice, 10),
	}
	if err := shareTokenData.Save(ctx, db); err != nil {
		return fmt.Errorf("[createShareTokenData] failed to save share token data: %w", err)
	}

	if err := updateCurrentSharePrice(ctx, db, strategy.VaultID, ev.SharePrice); err != nil {
		return fmt.Errorf("[createShareTokenData] failed to update current share price: %w", err)
	}
	return nil
}

func updateCurrentSharePrice(ctx context.Context, db *gorm.DB, vaultId string, sharePrice uint64) error {
	log.Infof("[updateCurrentSharePrice] Updating current share price...")
	vault := modelsss.Vault{ID: vaultId}
	ok, err := vault.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[updateCurrentSharePrice] failed to load vault: %w", err)
	}
	vault.CurrentSharePrice = strconv.FormatUint(sharePrice, 10)
	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[updateCurrentSharePrice] failed to save vault: %w", err)
	}

	token := modelsss.Token{ID: vault.ShareTokenID}
	ok, err = token.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[updateCurrentSharePrice] failed to load token: %w", err)
	}
	token.CurrentPrice = sharePrice
	if err := token.Save(ctx, db); err != nil {
		return fmt.Errorf("[updateCurrentSharePrice] failed to save token: %w", err)
	}

	return nil
}

func createReportResult(ctx context.Context, db *gorm.DB, event coremodel.Event, previousReport modelsss.StrategyReport, currentReport modelsss.StrategyReport) error {
	log.Infof("[createReportResult] Creating report result (latest vs current report)...")
	if currentReport.ID == previousReport.ID {
		log.Warnf("[createReportResult] Report result NOT created. Current report is the same as latest report")
		return nil
	}

	id := utils.GenerateId(event)
	strategyReportResult := modelsss.StrategyReportResult{ID: id}
	if _, err := strategyReportResult.Load(ctx, db); err != nil {
		return fmt.Errorf("[createReportResult] failed to load strategy report result: %w", err)
	}

	strategyReportResult.TransactionHash = event.TransactionSignature
	strategyReportResult.Timestamp = strconv.FormatInt(event.BlockTime, 10)
	strategyReportResult.BlockNumber = strconv.FormatUint(event.Slot, 10)
	strategyReportResult.PreviousReportID = previousReport.ID
	strategyReportResult.CurrentReportID = currentReport.ID
	strategyReportResult.StartTimestamp = previousReport.Timestamp
	strategyReportResult.EndTimestamp = currentReport.Timestamp

	currentReportTimestamp, err := utils.ParseBigIntFromString(currentReport.Timestamp)
	if err != nil {
		return fmt.Errorf("[createReportResult] failed to parse current report timestamp: %w", err)
	}
	previousReportTimestamp, err := utils.ParseBigIntFromString(previousReport.Timestamp)
	if err != nil {
		return fmt.Errorf("[createReportResult] failed to parse previous report timestamp: %w", err)
	}
	duration := new(big.Int).Sub(currentReportTimestamp, previousReportTimestamp)

	strategyReportResult.Duration = duration.String()

	////TODO: Check if this is correct
	//profit := currentReport.Gain - currentReport.Loss
	//msInDays := utils.MillisToDays(strategyReportResult.Duration)
	//
	//log.Infof("[createReportResult] Report Result - Start / End: %s / %s - Duration: %s (days %f) - Profit: %d",
	//	time.Unix(strategyReportResult.StartTimestamp, 0),
	//	time.Unix(strategyReportResult.EndTimestamp, 0),
	//	time.Unix(strategyReportResult.Duration, 0),
	//	msInDays, profit)
	//
	//if previousReport.CurrentDebt != 0 && msInDays != 0 {
	//	profitOverTotalDebt := profit / previousReport.CurrentDebt
	//	strategyReportResult.DurationPr = profitOverTotalDebt
	//}

	return nil
}
