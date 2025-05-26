package vault

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/monitoring"
	"github.com/Tsisar/solana-indexer/internal/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/internal/subgraph/events"
	"github.com/Tsisar/solana-indexer/internal/subgraph/library/account"
	"github.com/Tsisar/solana-indexer/internal/subgraph/library/accountant"
	"github.com/Tsisar/solana-indexer/internal/subgraph/library/report"
	"github.com/Tsisar/solana-indexer/internal/subgraph/library/token"
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"github.com/Tsisar/solana-indexer/internal/utils"
	"math/big"

	"gorm.io/gorm"
)

func Init(ctx context.Context, db *gorm.DB, ev events.VaultInitEvent, transaction events.Transaction) error {
	var err error
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err = vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}

	vault.Token, err = token.UpsertUnderlyingToken(ctx, db, ev)
	if err != nil {
		return fmt.Errorf("[vault] failed to get or create underlying token entity: %v", err)
	}

	vault.ShareToken, err = token.UpsertShareToken(ctx, db, ev)
	if err != nil {
		return fmt.Errorf("[vault] failed to get or create share token entity: %v", err)
	}

	acc, err := accountant.Init(ctx, db, ev.Accountant.String())
	if err != nil {
		return fmt.Errorf("[vault] failed to upsert accountant: %w", err)
	}

	vault.DepositLimit = ev.DepositLimit
	vault.Shutdown = false
	vault.Activation = transaction.Timestamp
	vault.AccountantID = acc.ID
	vault.MinUserDeposit = ev.MinUserDeposit
	vault.KycVerifiedOnly = ev.KycVerifiedOnly
	vault.DirectDepositEnabled = ev.DirectDepositEnabled
	vault.WhitelistedOnly = ev.WhitelistedOnly
	vault.ProfitMaxUnlockTime = ev.ProfitMaxUnlockTime
	vault.LastUpdate = transaction.Timestamp
	vault.MinTotalIdle = ev.MinimumTotalIdle
	vault.DirectWithdrawEnabled = ev.DirectWithdrawEnabled
	vault.UserDepositLimit = ev.UserDepositLimit

	if err = vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	return nil
}

func AddStrategy(ctx context.Context, db *gorm.DB, ev events.VaultAddStrategyEvent, transaction events.Transaction) error {
	strategy := subgraph.Strategy{
		ID: ev.StrategyKey.String(),
	}
	ok, err := strategy.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[vault] failed to load strategy: %w", err)
	}

	strategy.MaxDebt = ev.MaxDebt
	strategy.CurrentDebt = ev.CurrentDebt
	strategy.VaultID = ev.VaultKey.String()
	strategy.Activation = transaction.Timestamp

	if err := strategy.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save strategy: %w", err)
	}
	return nil
}

func Deposit(ctx context.Context, db *gorm.DB, ev events.VaultDepositEvent, transaction events.Transaction) error {
	// return ev.Authority.String()
	if err := account.UpdateAccount(ctx, db, ev.Authority.String(), ev.TokenAccount.String(), ev.ShareAccount.String()); err != nil {
		return fmt.Errorf("[vault] failed to update account: %w", err)
	}

	id := utils.GenerateId(transaction.Signature, transaction.EventIndex.String())
	deposit := subgraph.Deposit{ID: id}
	if _, err := deposit.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load deposit: %w", err)
	}

	deposit.Timestamp = ev.Timestamp
	deposit.BlockNumber = transaction.Slot
	deposit.AccountID = ev.Authority.String()
	deposit.VaultID = ev.VaultKey.String()
	deposit.TokenAmount = ev.Amount
	deposit.SharesMinted = ev.Share
	deposit.ShareTokenID = ev.ShareMint.String()
	deposit.TokenID = ev.TokenMint.String()
	deposit.SharePrice = ev.SharePrice
	if err := deposit.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save deposit: %w", err)
	}
	monitoring.Deposit(ctx, db, deposit)

	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err := vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	vault.TotalDebt = ev.TotalDebt
	vault.TotalIdle = ev.TotalIdle
	vault.TotalShare = ev.TotalShare
	vault.BalanceTokensIdle = ev.TotalIdle
	vault.BalanceTokens = ev.TotalShare
	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	if err := vaultPositionDeposit(ctx, db, ev); err != nil {
		return fmt.Errorf("[vault] failed to save vault position deposit: %w", err)
	}

	if err := UpdateCurrentSharePrice(ctx, db, vault.ID, ev.SharePrice); err != nil {
		return fmt.Errorf("[vault] failed to update current share price: %w", err)
	}

	return nil
}

func StrategyReported(ctx context.Context, db *gorm.DB, ev events.StrategyReportedEvent, transaction events.Transaction) error {
	if err := report.CreateReport(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[vault] failed to create report: %w", err)
	}

	if err := report.CreateReportEvent(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[vault] failed to create report event: %w", err)
	}

	if err := report.CreateShareTokenData(ctx, db, ev, transaction); err != nil {
		return fmt.Errorf("[vault] failed to create share token data: %w", err)
	}

	if err := UpdateCurrentSharePrice(ctx, db, ev.VaultKey.String(), ev.SharePrice); err != nil {
		return fmt.Errorf("[vault] failed to update current share price: %w", err)
	}
	return nil
}

func vaultPositionDeposit(ctx context.Context, db *gorm.DB, ev events.VaultDepositEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err := vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}

	t := subgraph.Token{ID: vault.TokenID}
	if _, err := t.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load token wallet: %w", err)
	}

	sum := ev.TotalDebt.Plus(&ev.TotalIdle)
	if sum == nil || sum.Int == nil {
		log.Warn("[vault] sum is nil, cannot compute price per share")
		return nil
	}
	pricePerShare := sum.Div(&ev.TotalShare) //(total_debt + total_idle) / total_shares

	product := ev.TotalShare.Mul(pricePerShare)
	if product == nil || product.Int == nil {
		log.Warn("[vault] product is nil, cannot compute balance position")
		return nil
	}
	balancePosition := product.Div(&t.Decimals) // total_shares * price_per_share

	id := utils.GenerateId(ev.VaultKey.String(), ev.Authority.String())
	accountVaultPosition := subgraph.AccountVaultPosition{ID: id}
	ok, err := accountVaultPosition.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load vault position: %w", err)
	}
	if ok {
		balanceTokens := accountVaultPosition.BalanceTokens.Plus(&ev.Amount)
		accountVaultPosition.BalanceTokens = *balanceTokens
		accountVaultPosition.BalanceShares = ev.TotalShare
	} else {
		accountVaultPosition.VaultID = vault.ID
		accountVaultPosition.AccountID = ev.Authority.String()
		accountVaultPosition.TokenID = vault.TokenID
		accountVaultPosition.ShareTokenID = vault.ShareTokenID
		accountVaultPosition.BalanceTokens = ev.Amount
		accountVaultPosition.BalanceShares = ev.TotalShare
		accountVaultPosition.BalancePosition = *balancePosition
	}

	if err := accountVaultPosition.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault position: %w", err)
	}
	return nil
}

func Withdraw(ctx context.Context, db *gorm.DB, ev events.VaultWithdrawlEvent, transaction events.Transaction) error {
	if err := account.UpdateAccount(ctx, db, ev.Authority.String(), ev.TokenAccount.String(), ev.ShareAccount.String()); err != nil {
		return fmt.Errorf("[vault] failed to update account: %w", err)
	}

	id := utils.GenerateId(transaction.Signature, transaction.EventIndex.String())
	withdrwal := subgraph.Withdrawal{ID: id}
	if _, err := withdrwal.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load withdrawal: %w", err)
	}
	withdrwal.Timestamp = ev.Timestamp
	withdrwal.BlockNumber = transaction.Slot
	withdrwal.AccountID = ev.Authority.String()
	withdrwal.VaultID = ev.VaultKey.String()
	withdrwal.TokenAmount = ev.AssetsToTransfer
	withdrwal.SharesBurnt = ev.SharesToBurn
	withdrwal.ShareTokenID = ev.ShareMint.String()
	withdrwal.TokenID = ev.TokenMint.String()
	withdrwal.SharePrice = ev.SharePrice
	if err := withdrwal.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save withdrawal: %w", err)
	}
	monitoring.Withdrawal(ctx, db, withdrwal)

	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err := vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}

	//TODO: This is as per fathom logic...
	vault.TotalIdle = ev.TotalIdle
	vault.TotalShare = ev.TotalShare
	vault.BalanceTokensIdle = ev.TotalIdle
	vault.BalanceTokens = ev.TotalShare

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	if err := vaultPositionWithdraw(ctx, db, ev); err != nil {
		return fmt.Errorf("[vault] failed to save vault position withdraw: %w", err)
	}

	if err := UpdateCurrentSharePrice(ctx, db, vault.ID, ev.SharePrice); err != nil {
		return fmt.Errorf("[vault] failed to update current share price: %w", err)
	}

	return nil
}

func vaultPositionWithdraw(ctx context.Context, db *gorm.DB, ev events.VaultWithdrawlEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err := vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}

	t := subgraph.Token{ID: vault.TokenID}
	if _, err := t.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load token wallet: %w", err)
	}

	//TODO: check if this is correct
	var pricePerShare *types.BigInt
	sum := vault.TotalDebt.Plus(&vault.TotalIdle)
	if vault.TotalShare.Int != nil && vault.TotalShare.Int.Sign() > 0 {
		pricePerShare = sum.Div(&vault.TotalShare)
	} else {
		log.Info("[mapping] totalShares is zero!")
		pricePerShare = sum
	}

	if t.Decimals.Int.Sign() == 0 {
		log.Warn("[mapping] token decimals is zero!")
		return nil
	}

	product := ev.TotalShare.Mul(pricePerShare)
	if product == nil || product.Int == nil {
		log.Warn("[vault] product is nil, cannot compute balance position")
		return nil
	}
	balancePosition := product.Div(&t.Decimals) // total_shares * price_per_share

	id := utils.GenerateId(ev.VaultKey.String(), ev.Authority.String())
	accountVaultPosition := subgraph.AccountVaultPosition{ID: id}
	ok, err := accountVaultPosition.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load vault position: %w", err)
	}
	if ok {
		accountVaultPosition.BalanceShares = vault.TotalShare
		balanceTokens := GetBalanceTokens(&accountVaultPosition.BalanceTokens, &ev.AssetsToTransfer)
		accountVaultPosition.BalanceTokens = *balanceTokens
		balanceProfit := GetBalanceProfit(
			&accountVaultPosition.BalanceShares,
			&accountVaultPosition.BalanceProfit,
			&accountVaultPosition.BalanceTokens,
			&ev.AssetsToTransfer)
		accountVaultPosition.BalanceProfit = *balanceProfit
		accountVaultPosition.BalancePosition = *balancePosition
		if err := accountVaultPosition.Save(ctx, db); err != nil {
			return fmt.Errorf("[vault] failed to save vault position: %w", err)
		}
	}

	return nil
}

func GetBalanceProfit(currentSharesBalance, currentProfit, currentAmount, withdrawAmount *types.BigInt) *types.BigInt {
	zero := types.ZeroBigInt()
	// Defensive: if any is nil, treat as zero
	if currentSharesBalance == nil || currentSharesBalance.Int == nil ||
		currentProfit == nil || currentProfit.Int == nil ||
		currentAmount == nil || currentAmount.Int == nil ||
		withdrawAmount == nil || withdrawAmount.Int == nil {
		return &zero
	}

	if currentSharesBalance.Int.Sign() == 0 {
		// User withdrawn all the shares
		switch withdrawAmount.Int.Cmp(currentAmount.Int) {
		case 1: // withdrawAmount > currentAmount → profit
			return currentProfit.Plus(&types.BigInt{
				Int: new(big.Int).Sub(withdrawAmount.Int, currentAmount.Int),
			})
		case -1: // withdrawAmount < currentAmount → loss
			return currentProfit.Sub(&types.BigInt{
				Int: new(big.Int).Sub(currentAmount.Int, withdrawAmount.Int),
			})
		default: // equal
			return currentProfit
		}
	}

	// User still has shares → return current profit
	return currentProfit
}

func GetBalanceTokens(current, withdraw *types.BigInt) *types.BigInt {
	zero := types.ZeroBigInt()
	if current == nil || current.Int == nil || withdraw == nil || withdraw.Int == nil {
		return &zero
	}

	if withdraw.Int.Cmp(current.Int) > 0 {
		return &zero
	}

	return current.Sub(withdraw)
}

func UpdateCurrentSharePrice(ctx context.Context, db *gorm.DB, vaultId string, sharePrice types.BigInt) error {
	log.Infof("[vault] Updating current share price...")
	vault := subgraph.Vault{ID: vaultId}
	ok, err := vault.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	vault.CurrentSharePrice = sharePrice
	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	t := subgraph.Token{ID: vault.ShareTokenID}
	ok, err = t.Load(ctx, db)
	if err != nil || !ok {
		return fmt.Errorf("[vault] failed to load token: %w", err)
	}
	t.CurrentPrice = sharePrice
	if err := t.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save token: %w", err)
	}
	monitoring.Token(t)

	return nil
}

func UpdateDepositLimit(ctx context.Context, db *gorm.DB, ev events.VaultUpdateDepositLimitEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err := vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	vault.DepositLimit = ev.NewLimit
	vault.LastUpdate = ev.Timestamp

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}
	return nil
}

func ShutDown(ctx context.Context, db *gorm.DB, ev events.VaultShutDownEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	if _, err := vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	vault.Shutdown = ev.Shutdown

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}
	return nil
}

func WithdrawalRequested(ctx context.Context, db *gorm.DB, ev events.WithdrawalRequestedEvent) error {
	id := utils.GenerateId(ev.User.String(), ev.Vault.String(), ev.Index.String())

	withdrawalRequest := subgraph.WithdrawalRequest{ID: id}
	if _, err := withdrawalRequest.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load withdrawal request: %w", err)
	}

	withdrawalRequest.User = ev.User.String()
	withdrawalRequest.VaultID = ev.Vault.String()
	withdrawalRequest.Recipient = ev.Recipient.String()
	withdrawalRequest.Shares = ev.Shares
	withdrawalRequest.Amount = ev.Amount
	withdrawalRequest.MaxLoss = ev.MaxLoss
	withdrawalRequest.FeeShares = ev.FeeShares
	withdrawalRequest.Index = ev.Index
	withdrawalRequest.Open = true
	withdrawalRequest.Status = "open"
	withdrawalRequest.Timestamp = ev.Timestamp
	withdrawalRequest.PriorityFees = &ev.PriorityFee

	if err := withdrawalRequest.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save withdrawal request: %w", err)
	}

	if err := updatePriorityFeeOnVault(ctx, db, ev.Vault.String(), &ev.PriorityFee, "open"); err != nil {
		return fmt.Errorf("[vault] failed to update priority fee on vault: %w", err)
	}

	return nil
}

func WithdrawalRequestFulfilled(ctx context.Context, db *gorm.DB, ev events.WithdrawalRequestFulfilledEvent) error {
	id := utils.GenerateId(ev.User.String(), ev.Vault.String(), ev.Index.String())

	withdrawalRequest := subgraph.WithdrawalRequest{ID: id}
	ok, err := withdrawalRequest.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load withdrawal request: %w", err)
	}
	if !ok {
		log.Warnf("[vault] withdrawal request not found: %s", id)
		return nil
	}

	withdrawalRequest.Status = "fulfilled"
	withdrawalRequest.Open = false
	withdrawalRequest.Timestamp = ev.Timestamp
	withdrawalRequest.Amount = ev.Amount

	if err := withdrawalRequest.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save withdrawal request: %w", err)
	}

	if err := updatePriorityFeeOnVault(ctx, db, ev.Vault.String(), withdrawalRequest.PriorityFees, "fulfilled"); err != nil {
		return fmt.Errorf("[vault] failed to update priority fee on vault: %w", err)
	}

	return nil
}

func WithdrawalRequestCanceled(ctx context.Context, db *gorm.DB, ev events.WithdrawalRequestCanceledEvent) error {
	id := utils.GenerateId(ev.User.String(), ev.Vault.String(), ev.Index.String())

	withdrawalRequest := subgraph.WithdrawalRequest{ID: id}
	ok, err := withdrawalRequest.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load withdrawal request: %w", err)
	}
	if !ok {
		log.Warnf("[vault] withdrawal request not found: %s", id)
		return nil
	}

	withdrawalRequest.Status = "canceled"
	withdrawalRequest.Open = false
	withdrawalRequest.Timestamp = ev.Timestamp

	if err := withdrawalRequest.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save withdrawal request: %w", err)
	}

	if err := updatePriorityFeeOnVault(ctx, db, ev.Vault.String(), withdrawalRequest.PriorityFees, "canceled"); err != nil {
		return fmt.Errorf("[vault] failed to update priority fee on vault: %w", err)
	}

	return nil
}

func updatePriorityFeeOnVault(ctx context.Context, db *gorm.DB, vaultID string, fee *types.BigInt, status string) error {
	vault := subgraph.Vault{ID: vaultID}
	if _, err := vault.Load(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}

	if status == "open" {
		vault.TotalPriorityFees = utils.Val(vault.TotalPriorityFees.Plus(fee))
	} else {
		vault.TotalPriorityFees = utils.Val(vault.TotalPriorityFees.Sub(fee))
	}
	return nil
}

func UpdateWhiteListOnly(ctx context.Context, db *gorm.DB, ev events.VaultUpdateWhitelistedOnlyEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	ok, err := vault.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	if !ok {
		log.Warnf("[vault] vault not found: %s", ev.VaultKey.String())
		return nil
	}

	vault.WhitelistedOnly = ev.NewWhitelistedOnly
	vault.LastUpdate = ev.Timestamp

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	return nil
}

func UpdateAccountant(ctx context.Context, db *gorm.DB, ev events.VaultUpdateAccountantEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	ok, err := vault.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	if !ok {
		log.Warnf("[vault] vault not found: %s", ev.VaultKey.String())
		return nil
	}

	vault.AccountantID = ev.NewAccountant.String()
	vault.LastUpdate = ev.Timestamp

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	return nil
}

func UpdateUserDepositLimit(ctx context.Context, db *gorm.DB, ev events.VaultUpdateUserDepositLimitEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	ok, err := vault.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	if !ok {
		log.Warnf("[vault] vault not found: %s", ev.VaultKey.String())
		return nil
	}

	vault.UserDepositLimit = ev.NewUserDepositLimit
	vault.LastUpdate = ev.Timestamp

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	return nil
}

func UpdateDirectWithdrawEnabled(ctx context.Context, db *gorm.DB, ev events.VaultUpdateDirectWithdrawEnabledEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	ok, err := vault.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	if !ok {
		log.Warnf("[vault] vault not found: %s", ev.VaultKey.String())
		return nil
	}

	vault.DirectWithdrawEnabled = ev.NewDirectWithdrawEnabled
	vault.LastUpdate = ev.Timestamp

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	return nil
}

func UpdateMinTotalIdle(ctx context.Context, db *gorm.DB, ev events.VaultUpdateMinTotalIdleEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	ok, err := vault.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	if !ok {
		log.Warnf("[vault] vault not found: %s", ev.VaultKey.String())
		return nil
	}

	vault.MinTotalIdle = ev.NewMinTotalIdle
	vault.LastUpdate = ev.Timestamp

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	return nil
}

func UpdateProfitMaxUnlockTime(ctx context.Context, db *gorm.DB, ev events.VaultUpdateProfitMaxUnlockTimeEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	ok, err := vault.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	if !ok {
		log.Warnf("[vault] vault not found: %s", ev.VaultKey.String())
		return nil
	}

	vault.ProfitMaxUnlockTime = ev.NewProfitMaxUnlockTime
	vault.LastUpdate = ev.Timestamp

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	return nil
}

func UpdateMinUserDeposit(ctx context.Context, db *gorm.DB, ev events.VaultUpdateMinUserDepositEvent) error {
	vault := subgraph.Vault{ID: ev.VaultKey.String()}
	ok, err := vault.Load(ctx, db)
	if err != nil {
		return fmt.Errorf("[vault] failed to load vault: %w", err)
	}
	if !ok {
		log.Warnf("[vault] vault not found: %s", ev.VaultKey.String())
		return nil
	}

	vault.MinUserDeposit = ev.NewMinUserDeposit
	vault.LastUpdate = ev.Timestamp

	if err := vault.Save(ctx, db); err != nil {
		return fmt.Errorf("[vault] failed to save vault: %w", err)
	}

	return nil
}
