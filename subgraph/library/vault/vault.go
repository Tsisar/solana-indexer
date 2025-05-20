package vault

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/subgraph/events"
	"github.com/Tsisar/solana-indexer/subgraph/library/account"
	"github.com/Tsisar/solana-indexer/subgraph/library/accountant"
	"github.com/Tsisar/solana-indexer/subgraph/library/report"
	"github.com/Tsisar/solana-indexer/subgraph/library/token"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"github.com/Tsisar/solana-indexer/subgraph/utils"

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

	if err := vaultPositionDeposit(ctx, db, ev, transaction); err != nil {
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

func vaultPositionDeposit(ctx context.Context, db *gorm.DB, ev events.VaultDepositEvent, transaction events.Transaction) error {
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

	return nil
}
