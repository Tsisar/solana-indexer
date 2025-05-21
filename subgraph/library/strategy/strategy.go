package strategy

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/subgraph/events"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"github.com/Tsisar/solana-indexer/subgraph/utils"
	"gorm.io/gorm"
	"math/big"
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
	freeFunds := subgraph.DeployFunds{ID: id}
	if _, err := freeFunds.Load(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to load deploy funds: %w", err)
	}
	freeFunds.StrategyID = ev.AccountKey.String()
	freeFunds.Amount = ev.Amount
	freeFunds.Timestamp = ev.Timestamp

	if err := freeFunds.Save(ctx, db); err != nil {
		return fmt.Errorf("[strategy] failed to save deploy funds: %w", err)
	}
	return nil
}

func AfterOrcaSwap(ctx context.Context, db *gorm.DB, ev events.OrcaAfterSwapEvent) error {
	totalAssets := ev.IdleUnderlying.Plus(&ev.TotalInvested)
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
	if vaultTotalAllocation == nil || vaultTotalAllocation.Sign() == 0 {
		strategy.TotalAllocationPercent = utils.Val(strategy.TotalAllocation.SafeDiv(vaultTotalAllocation).Mul(&hundred))
	}

	if ev.Buy {
		strategy.EffectiveInvestedAmount = utils.Val(strategy.EffectiveInvestedAmount.Plus(&ev.Amount))
	} else {
		totalAssets = &ev.TotalInvested
		tokensSold := ev.AssetBalanceBefore.Sub(&ev.AssetBalanceAfter)
		costBasisForSale := &types.BigInt{Int: big.NewInt(0)}

		if ev.AssetBalanceAfter.Sign() != 0 {
			mul := strategy.EffectiveInvestedAmount.Mul(tokensSold)
			costBasisForSale = mul.Div(&ev.AssetBalanceBefore)
		}

		strategy.EffectiveInvestedAmount = utils.Val(strategy.EffectiveInvestedAmount.Sub(costBasisForSale))
	}
	strategy.ProfitOrLoss = utils.Val(totalAssets.Sub(&strategy.EffectiveInvestedAmount).ToBigDecimal())

	if strategy.EffectiveInvestedAmount.Sign() != 0 {
		strategy.ProfitOrLossPercent = utils.Val(strategy.ProfitOrLoss.SafeDiv(strategy.EffectiveInvestedAmount.ToBigDecimal()).Mul(&hundred))
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
