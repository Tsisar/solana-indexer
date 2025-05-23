package events

// Code generated by generate_events.go; DO NOT EDIT.

import "github.com/gagliardetto/solana-go"

// DepositLimitSetEvent event struct
type DepositLimitSetEvent struct {
	AccountKey   solana.PublicKey `borsh:"account_key"`
	DepositLimit uint64           `borsh:"deposit_limit"`
	Timestamp    int64            `borsh:"timestamp"`
}

// EmergencyWithdrawEvent event struct
type EmergencyWithdrawEvent struct {
	StrategyKey      solana.PublicKey `borsh:"strategy_key"`
	VaultKey         solana.PublicKey `borsh:"vault_key"`
	AssetMint        solana.PublicKey `borsh:"asset_mint"`
	Recipient        solana.PublicKey `borsh:"recipient"`
	RedeemableAmount uint64           `borsh:"redeemable_amount"`
	ScaledRatio      uint64           `borsh:"scaled_ratio"`
	Timestamp        int64            `borsh:"timestamp"`
}

// FundManagerDeployFundsEvent event struct
type FundManagerDeployFundsEvent struct {
	AccountKey     solana.PublicKey `borsh:"account_key"`
	Vault          solana.PublicKey `borsh:"vault"`
	Amount         uint64           `borsh:"amount"`
	DeployedAmount uint64           `borsh:"deployed_amount"`
	TotalInvested  uint64           `borsh:"total_invested"`
	TotalDeployed  uint64           `borsh:"total_deployed"`
	Timestamp      int64            `borsh:"timestamp"`
}

// FundManagerEmergencyWithdrawEvent event struct
type FundManagerEmergencyWithdrawEvent struct {
	AccountKey        solana.PublicKey `borsh:"account_key"`
	Vault             solana.PublicKey `borsh:"vault"`
	Amount            uint64           `borsh:"amount"`
	AmountTransferred uint64           `borsh:"amount_transferred"`
	TotalInvested     uint64           `borsh:"total_invested"`
	Timestamp         int64            `borsh:"timestamp"`
}

// FundManagerFreeFundsEvent event struct
type FundManagerFreeFundsEvent struct {
	AccountKey        solana.PublicKey `borsh:"account_key"`
	Vault             solana.PublicKey `borsh:"vault"`
	Amount            uint64           `borsh:"amount"`
	AmountTransferred uint64           `borsh:"amount_transferred"`
	TotalInvested     uint64           `borsh:"total_invested"`
	TotalFreed        uint64           `borsh:"total_freed"`
	Timestamp         int64            `borsh:"timestamp"`
}

// FundManagerHarvestAndReportEvent event struct
type FundManagerHarvestAndReportEvent struct {
	AccountKey    solana.PublicKey `borsh:"account_key"`
	Vault         solana.PublicKey `borsh:"vault"`
	TotalInvested uint64           `borsh:"total_invested"`
	TotalAssets   uint64           `borsh:"total_assets"`
	Timestamp     int64            `borsh:"timestamp"`
}

// FundManagerStrategyStateUpdateEvent event struct
type FundManagerStrategyStateUpdateEvent struct {
	AccountKey    solana.PublicKey `borsh:"account_key"`
	Vault         solana.PublicKey `borsh:"vault"`
	TotalAssets   uint64           `borsh:"total_assets"`
	TotalInvested uint64           `borsh:"total_invested"`
	TotalIdle     uint64           `borsh:"total_idle"`
	TotalDeployed uint64           `borsh:"total_deployed"`
	TotalFreed    uint64           `borsh:"total_freed"`
	Timestamp     int64            `borsh:"timestamp"`
}

// HarvestAndReportDTFEvent event struct
type HarvestAndReportDTFEvent struct {
	AccountKey  solana.PublicKey `borsh:"account_key"`
	TotalAssets uint64           `borsh:"total_assets"`
	Timestamp   int64            `borsh:"timestamp"`
}

// MinDeployAmountSetEvent event struct
type MinDeployAmountSetEvent struct {
	AccountKey      solana.PublicKey `borsh:"account_key"`
	MinDeployAmount uint64           `borsh:"min_deploy_amount"`
	Timestamp       int64            `borsh:"timestamp"`
}

// OrcaAfterSwapEvent event struct
type OrcaAfterSwapEvent struct {
	AccountKey              solana.PublicKey `borsh:"account_key"`
	Vault                   solana.PublicKey `borsh:"vault"`
	Buy                     bool             `borsh:"buy"`
	Amount                  uint64           `borsh:"amount"`
	TotalInvested           uint64           `borsh:"total_invested"`
	WhirlpoolId             solana.PublicKey `borsh:"whirlpool_id"`
	UnderlyingMint          solana.PublicKey `borsh:"underlying_mint"`
	UnderlyingDecimals      uint8            `borsh:"underlying_decimals"`
	AssetMint               solana.PublicKey `borsh:"asset_mint"`
	AssetAmount             uint64           `borsh:"asset_amount"`
	AssetDecimals           uint8            `borsh:"asset_decimals"`
	TotalAssets             uint64           `borsh:"total_assets"`
	IdleUnderlying          uint64           `borsh:"idle_underlying"`
	AToBForPurchase         bool             `borsh:"a_to_b_for_purchase"`
	UnderlyingBalanceBefore uint64           `borsh:"underlying_balance_before"`
	UnderlyingBalanceAfter  uint64           `borsh:"underlying_balance_after"`
	AssetBalanceBefore      uint64           `borsh:"asset_balance_before"`
	AssetBalanceAfter       uint64           `borsh:"asset_balance_after"`
	Timestamp               int64            `borsh:"timestamp"`
}

// OrcaInitEvent event struct
type OrcaInitEvent struct {
	AccountKey      solana.PublicKey `borsh:"account_key"`
	WhirlpoolId     solana.PublicKey `borsh:"whirlpool_id"`
	AssetMint       solana.PublicKey `borsh:"asset_mint"`
	AssetDecimals   uint8            `borsh:"asset_decimals"`
	AToBForPurchase bool             `borsh:"a_to_b_for_purchase"`
}

// SetPerformanceFeeEvent event struct
type SetPerformanceFeeEvent struct {
	AccountKey solana.PublicKey `borsh:"account_key"`
	Fee        uint64           `borsh:"fee"`
}

// StrategyDeployFundsEvent event struct
type StrategyDeployFundsEvent struct {
	AccountKey solana.PublicKey `borsh:"account_key"`
	Amount     uint64           `borsh:"amount"`
	Timestamp  int64            `borsh:"timestamp"`
}

// StrategyDepositEvent event struct
type StrategyDepositEvent struct {
	AccountKey  solana.PublicKey `borsh:"account_key"`
	Amount      uint64           `borsh:"amount"`
	TotalAssets uint64           `borsh:"total_assets"`
}

// StrategyFreeFundsEvent event struct
type StrategyFreeFundsEvent struct {
	AccountKey solana.PublicKey `borsh:"account_key"`
	Amount     uint64           `borsh:"amount"`
	Timestamp  int64            `borsh:"timestamp"`
}

// StrategyInitEvent event struct
type StrategyInitEvent struct {
	AccountKey         solana.PublicKey `borsh:"account_key"`
	StrategyType       string           `borsh:"strategy_type"`
	Vault              solana.PublicKey `borsh:"vault"`
	UnderlyingMint     solana.PublicKey `borsh:"underlying_mint"`
	UnderlyingTokenAcc solana.PublicKey `borsh:"underlying_token_acc"`
	UnderlyingDecimals uint8            `borsh:"underlying_decimals"`
	DepositLimit       uint64           `borsh:"deposit_limit"`
	DepositPeriodEnds  int64            `borsh:"deposit_period_ends"`
	LockPeriodEnds     int64            `borsh:"lock_period_ends"`
}

// StrategyReallocEvent event struct
type StrategyReallocEvent struct {
	Strategy  solana.PublicKey `borsh:"strategy"`
	NewSize   uint64           `borsh:"new_size"`
	Timestamp int64            `borsh:"timestamp"`
}

// StrategyShutdownEvent event struct
type StrategyShutdownEvent struct {
	AccountKey solana.PublicKey `borsh:"account_key"`
	Shutdown   bool             `borsh:"shutdown"`
	Timestamp  int64            `borsh:"timestamp"`
}

// StrategyWithdrawEvent event struct
type StrategyWithdrawEvent struct {
	AccountKey  solana.PublicKey `borsh:"account_key"`
	Amount      uint64           `borsh:"amount"`
	TotalAssets uint64           `borsh:"total_assets"`
}

// TotalInvestedUpdatedEvent event struct
type TotalInvestedUpdatedEvent struct {
	AccountKey            solana.PublicKey `borsh:"account_key"`
	Vault                 solana.PublicKey `borsh:"vault"`
	PreviousTotalInvested uint64           `borsh:"previous_total_invested"`
	TotalInvested         uint64           `borsh:"total_invested"`
	Timestamp             int64            `borsh:"timestamp"`
}
