package strategy

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/monitoring"
	"github.com/Tsisar/solana-indexer/internal/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/internal/subgraph/events"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"github.com/Tsisar/solana-indexer/internal/utils"
	"gorm.io/gorm"
)

func Init(ctx context.Context, db *gorm.DB, ev events.StrategyInitEvent) error {
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}

	strategy.StrategyType = ev.StrategyType
	strategy.DepositLimit = ev.DepositLimit
	strategy.DepositPeriodEnds = ev.DepositPeriodEnds
	strategy.LockPeriodEnds = ev.LockPeriodEnds
	strategy.VaultID = ev.Vault.String()
	strategy.UnderlyingMint = ev.UnderlyingMint.String()
	strategy.UnderlyingTokenAcc = ev.UnderlyingTokenAcc.String()
	strategy.UnderlyingDecimals = ev.UnderlyingDecimals

	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}
	return nil
}

func Deposit(ctx context.Context, db *gorm.DB, ev events.StrategyDepositEvent) error {
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	strategy.TotalAssets = ev.TotalAssets
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}

	return nil
}

func Withdraw(ctx context.Context, db *gorm.DB, ev events.StrategyWithdrawEvent) error {
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	strategy.TotalAssets = ev.TotalAssets
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}

	return nil
}

func UpdateCurrentDebt(ctx context.Context, db *gorm.DB, ev events.UpdatedCurrentDebtForStrategyEvent) error {
	strategy := subgraph.Strategy{ID: ev.StrategyKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	strategy.CurrentDebt = ev.NewDebt
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}

	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err := vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load vault: %w", err)
	}
	vault.TotalDebt = ev.TotalDebt
	vault.TotalIdle = ev.TotalIdle
	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save vault: %w", err)
	}

	return nil
}

func UpdatePerformanceFee(ctx context.Context, db *gorm.DB, ev events.SetPerformanceFeeEvent) error {
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	strategy.PerformanceFees = ev.Fee
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}

	return nil
}

func UpdateDTFReport(ctx context.Context, db *gorm.DB, ev events.HarvestAndReportDTFEvent) error {
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	if _, err := strategy.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}

	id := utils.GenerateId(ev.AccountKey.String(), ev.Timestamp.String())
	dtfReport := subgraph.DTFReport{ID: id}
	if _, err := dtfReport.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load dtf report: %w", err)
	}
	dtfReport.TotalAssets = ev.TotalAssets
	dtfReport.Timestamp = ev.Timestamp

	if err := dtfReport.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save dtf report: %w", err)
	}

	strategy.DtfReportID = &id
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}

	return nil
}

func DeployFunds(ctx context.Context, db *gorm.DB, ev events.StrategyDeployFundsEvent) error {
	id := utils.GenerateId(ev.AccountKey.String(), ev.Timestamp.String())
	deployFunds := subgraph.DeployFunds{ID: id}
	if _, err := deployFunds.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load deploy funds: %w", err)
	}
	deployFunds.StrategyID = ev.AccountKey.String()
	deployFunds.Amount = ev.Amount
	deployFunds.Timestamp = ev.Timestamp

	if err := deployFunds.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save deploy funds: %w", err)
	}
	return nil
}

func FreeFunds(ctx context.Context, db *gorm.DB, ev events.StrategyFreeFundsEvent) error {
	id := utils.GenerateId(ev.AccountKey.String(), ev.Timestamp.String())
	freeFunds := subgraph.FreeFunds{ID: id}
	if _, err := freeFunds.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load free funds: %w", err)
	}
	freeFunds.StrategyID = ev.AccountKey.String()
	freeFunds.Amount = ev.Amount
	freeFunds.Timestamp = ev.Timestamp

	if err := freeFunds.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save free funds: %w", err)
	}

	// Update strategy to track total freed funds for PnL calculations
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	if ok, err := strategy.Load(ctx, db); err == nil && ok {
		// This could be used to track cumulative freed funds if needed
		// For now, the main fix is in AfterOrcaSwap function
		log.Debugf("[strategy] FreeFunds event processed for strategy %s, amount: %s", 
			strategy.ID, ev.Amount.String())
	}

	return nil
}

// validateWithdrawalRequestForPnL checks if a withdrawal request should be considered for PnL adjustment
func validateWithdrawalRequestForPnL(request *subgraph.WithdrawalRequest, currentEventTime types.BigInt) bool {
	if request == nil {
		return false
	}
	
	// Must be truly open/pending
	if !request.Open || request.Status != "open" {
		return false
	}
	
	// Additional safety: only count recent withdrawal requests to avoid stale data
	// Requests older than 1 hour might be stale due to processing delays
	requestAge := currentEventTime.Sub(&request.Timestamp).Int64()
	maxAgeSeconds := int64(3600) // 1 hour in seconds
	
	if requestAge > maxAgeSeconds {
		log.Warnf("[strategy] Ignoring old withdrawal request (age: %d seconds): %s", 
			requestAge, request.ID)
		return false
	}
	
	// Must have a valid amount
	if request.Amount.Sign() <= 0 {
		return false
	}
	
	return true
}

func AfterOrcaSwap(ctx context.Context, db *gorm.DB, ev events.OrcaAfterSwapEvent) error {
	totalAssets := &ev.TotalInvested
	vaultTotalAllocation, err := getTotalAllocationAfterAfterOrcaSwap(ctx, db, ev, totalAssets)
	if err != nil {
		return fmt.Errorf("[strategy] failed to get total allocation after orca swap: %w", err)
	}

	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	ok, err := strategy.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	if !ok {
		log.Warnf("[strategy] strategy not found: %s", ev.AccountKey.String())
		return nil
	}
	hundred := types.NewBigDecimalFromFloat(100)

	strategy.TotalAllocation = utils.Val(totalAssets.ToBigDecimal())
	log.Debugf("[strategy] total allocation: %s", strategy.TotalAllocation.String())

	if vaultTotalAllocation != nil && vaultTotalAllocation.Sign() != 0 {
		strategy.TotalAllocationPercent = utils.Val(strategy.TotalAllocation.SafeDiv(vaultTotalAllocation).Mul(&hundred))
		log.Debugf("[strategy] total allocation percent: %s", strategy.TotalAllocationPercent.String())
	}

	// --- PnL Calculations ---
	// Compare total assets (current value) with current debt (allocated capital)
	// For sell swaps (free_funds), we need to adjust currentDebt to account for pending withdrawals
	
	// Convert raw BigInt values to BigDecimal for calculation
	totalAssetsBD := totalAssets.ToBigDecimal()
	currentDebtBD := strategy.CurrentDebt.ToBigDecimal()
	
	// Check if this is a sell swap (free_funds operation)
	adjustedCurrentDebt := currentDebtBD
	if !ev.Buy {
		// This is a sell swap (free_funds), calculate pending withdrawal amount
		// Get total pending withdrawal requests for this vault to estimate debt reduction
		vault := subgraph.Vault{ID: ev.Vault.String()}
		if _, err := vault.Load(ctx, db); err == nil {
			// Calculate total pending withdrawal amount with more robust filtering
			var totalPendingWithdrawals types.BigInt
			totalPendingWithdrawals.Zero()
			pendingCount := 0
			
			for _, request := range vault.WithdrawalRequests {
				if validateWithdrawalRequestForPnL(request, ev.Timestamp) {
					totalPendingWithdrawals = *totalPendingWithdrawals.Plus(&request.Amount)
					pendingCount++
				}
			}
			
			// Only adjust currentDebt if we have valid pending withdrawals
			if totalPendingWithdrawals.Sign() > 0 {
				pendingWithdrawalsBD := totalPendingWithdrawals.ToBigDecimal()
				
				// Cap the adjustment to prevent over-correction
				// Don't reduce debt by more than the current debt amount
				maxReduction := currentDebtBD
				if pendingWithdrawalsBD.Sub(maxReduction).Sign() > 0 {
					pendingWithdrawalsBD = maxReduction
					log.Warnf("[strategy] Capping pending withdrawal adjustment to currentDebt: %s", 
						maxReduction.String())
				}
				
				adjustedCurrentDebt = currentDebtBD.Sub(pendingWithdrawalsBD)
				
				// Ensure adjusted debt doesn't go negative
				if adjustedCurrentDebt.Sign() < 0 {
					adjustedCurrentDebt = new(types.BigDecimal)
				}
				
				log.Debugf("[strategy] Adjusted currentDebt from %s to %s (pending withdrawals: %s, count: %d)", 
					currentDebtBD.String(), adjustedCurrentDebt.String(), pendingWithdrawalsBD.String(), pendingCount)
			} else {
				log.Debugf("[strategy] No valid pending withdrawals found for debt adjustment")
			}
		} else {
			log.Warnf("[strategy] Failed to load vault for pending withdrawal calculation: %v", err)
		}
	}

	// Absolute PnL Calculation using adjusted debt
	profitOrLoss := totalAssetsBD.Sub(adjustedCurrentDebt)
	strategy.ProfitOrLoss = utils.Val(profitOrLoss)
	log.Debugf("[strategy] PnL: %s (totalAssets: %s, adjustedDebt: %s)", 
		strategy.ProfitOrLoss.String(), totalAssetsBD.String(), adjustedCurrentDebt.String())

	// PnL in Percentage (%) Calculation
	if adjustedCurrentDebt.Sign() != 0 {
		profitOrLossPercent := strategy.ProfitOrLoss.SafeDiv(adjustedCurrentDebt).Mul(&hundred)
		strategy.ProfitOrLossPercent = utils.Val(profitOrLossPercent)
		log.Debugf("[strategy] PnL (Percent): %s", strategy.ProfitOrLossPercent.String())
	} else {
		strategy.ProfitOrLossPercent.Zero()
		log.Debugf("[strategy] adjusted debt is zero, PnL percent is zero")
	}

	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}

	return nil
}

func getTotalAllocationAfterAfterOrcaSwap(ctx context.Context, db *gorm.DB, ev events.OrcaAfterSwapEvent, totalAssets *types.BigInt) (*types.BigDecimal, error) {
	zero := types.ZeroBigDecimal()
	vault := subgraph.Vault{ID: ev.Vault.String()}
	ok, err := vault.Load(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("[strategy] failed to load vault: %w", err)
	}
	if !ok {
		log.Debugf("[strategy] vault not found: %s", ev.Vault.String())
		return &zero, nil
	}

	totalAssetsDecimal := totalAssets.ToBigDecimal()
	totalAllocation := &zero
	strategies := vault.Strategies
	for _, strategy := range strategies {
		if strategy != nil && strategy.ID != ev.AccountKey.String() {
			totalAllocation = totalAllocation.Plus(&strategy.TotalAllocation)
		}
	}
	vault.TotalAllocation = utils.Val(totalAllocation.Plus(totalAssetsDecimal))

	if err := vault.Save(ctx, db); err != nil {
		return nil, fmt.Errorf("[strategy] failed to save vault: %w", err)
	}
	return &vault.TotalAllocation, nil
}

func InitOrca(ctx context.Context, db *gorm.DB, ev events.OrcaInitEvent) error {
	strategy := subgraph.Strategy{ID: ev.AccountKey.String()}
	ok, err := strategy.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	if !ok {
		log.Warnf("[strategy] strategy not found: %s", ev.AccountKey.String())
		return nil
	}

	tokenMint := subgraph.Token{ID: ev.AssetMint.String()}
	if _, err := tokenMint.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load token: %w", err)
	}
	tokenMint.Decimals = ev.AssetDecimals
	if err := tokenMint.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save token: %w", err)
	}
	monitoring.Token(tokenMint)

	strategy.AssetID = &tokenMint.ID
	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}
	return nil
}

func Remove(ctx context.Context, db *gorm.DB, ev events.VaultRemoveStrategyEvent) error {
	strategy := subgraph.Strategy{ID: ev.StrategyKey.String()}
	ok, err := strategy.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[strategy] failed to load strategy: %w", err)
	}
	if !ok {
		log.Warnf("[strategy] strategy not found: %s", ev.StrategyKey.String())
		return nil
	}
	strategy.Removed = true
	strategy.RemovedTimestamp = ev.RemovedAt

	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save strategy: %w", err)
	}
	return nil

}
