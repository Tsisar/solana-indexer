package events

// Code generated by generate_events.go; DO NOT EDIT.

import "github.com/gagliardetto/solana-go"

// StrategyReportedEvent event struct
type StrategyReportedEvent struct {
	VaultKey     solana.PublicKey `borsh:"vault_key"`
	StrategyKey  solana.PublicKey `borsh:"strategy_key"`
	Gain         uint64           `borsh:"gain"`
	Loss         uint64           `borsh:"loss"`
	CurrentDebt  uint64           `borsh:"current_debt"`
	ProtocolFees uint64           `borsh:"protocol_fees"`
	TotalFees    uint64           `borsh:"total_fees"`
	TotalShares  uint64           `borsh:"total_shares"`
	SharePrice   uint64           `borsh:"share_price"`
	Timestamp    int64            `borsh:"timestamp"`
}

// UpdatedCurrentDebtForStrategyEvent event struct
type UpdatedCurrentDebtForStrategyEvent struct {
	VaultKey    solana.PublicKey `borsh:"vault_key"`
	StrategyKey solana.PublicKey `borsh:"strategy_key"`
	TotalIdle   uint64           `borsh:"total_idle"`
	TotalDebt   uint64           `borsh:"total_debt"`
	NewDebt     uint64           `borsh:"new_debt"`
}

// VaultAddStrategyEvent event struct
type VaultAddStrategyEvent struct {
	VaultKey    solana.PublicKey `borsh:"vault_key"`
	StrategyKey solana.PublicKey `borsh:"strategy_key"`
	CurrentDebt uint64           `borsh:"current_debt"`
	MaxDebt     uint64           `borsh:"max_debt"`
	LastUpdate  int64            `borsh:"last_update"`
	IsActive    bool             `borsh:"is_active"`
}

// VaultDepositEvent event struct
type VaultDepositEvent struct {
	VaultKey     solana.PublicKey `borsh:"vault_key"`
	TotalDebt    uint64           `borsh:"total_debt"`
	TotalIdle    uint64           `borsh:"total_idle"`
	TotalShare   uint64           `borsh:"total_share"`
	Amount       uint64           `borsh:"amount"`
	Share        uint64           `borsh:"share"`
	TokenAccount solana.PublicKey `borsh:"token_account"`
	ShareAccount solana.PublicKey `borsh:"share_account"`
	TokenMint    solana.PublicKey `borsh:"token_mint"`
	ShareMint    solana.PublicKey `borsh:"share_mint"`
	Authority    solana.PublicKey `borsh:"authority"`
	SharePrice   uint64           `borsh:"share_price"`
	Timestamp    int64            `borsh:"timestamp"`
}

// VaultEmergencyWithdrawEvent event struct
type VaultEmergencyWithdrawEvent struct {
	VaultKey            solana.PublicKey `borsh:"vault_key"`
	Recipient           solana.PublicKey `borsh:"recipient"`
	Shares              uint64           `borsh:"shares"`
	VaultTotalShares    uint64           `borsh:"vault_total_shares"`
	ScaledRatio         uint64           `borsh:"scaled_ratio"`
	StrategiesProcessed uint64           `borsh:"strategies_processed"`
	Timestamp           int64            `borsh:"timestamp"`
}

// VaultInitEvent event struct
type VaultInitEvent struct {
	VaultKey        solana.PublicKey `borsh:"vault_key"`
	UnderlyingToken TokenData        `borsh:"underlying_token"`
	Accountant      solana.PublicKey `borsh:"accountant"`
	ShareToken      TokenData        `borsh:"share_token"`
	DepositLimit    uint64           `borsh:"deposit_limit"`
	UserDepositLimit      uint64           `borsh:"user_deposit_limit"`
	MinUserDeposit        uint64           `borsh:"min_user_deposit"`
	KycVerifiedOnly       bool             `borsh:"kyc_verified_only"`
	DirectDepositEnabled  bool             `borsh:"direct_deposit_enabled"`
	DirectWithdrawEnabled bool             `borsh:"direct_withdraw_enabled"`
	MinimumTotalIdle      uint64           `borsh:"minimum_total_idle"`
	WhitelistedOnly       bool             `borsh:"whitelisted_only"`
	ProfitMaxUnlockTime   uint64           `borsh:"profit_max_unlock_time"`
}

// VaultRemoveStrategyEvent event struct
type VaultRemoveStrategyEvent struct {
	VaultKey    solana.PublicKey `borsh:"vault_key"`
	StrategyKey solana.PublicKey `borsh:"strategy_key"`
	RemovedAt   int64            `borsh:"removed_at"`
}

// VaultShutDownEvent event struct
type VaultShutDownEvent struct {
	VaultKey solana.PublicKey `borsh:"vault_key"`
	Shutdown bool             `borsh:"shutdown"`
}

// VaultUpdateAccountantEvent event struct
type VaultUpdateAccountantEvent struct {
	VaultKey      solana.PublicKey `borsh:"vault_key"`
	NewAccountant solana.PublicKey `borsh:"new_accountant"`
	Timestamp     int64            `borsh:"timestamp"`
}

// VaultUpdateDepositLimitEvent event struct
type VaultUpdateDepositLimitEvent struct {
	VaultKey  solana.PublicKey `borsh:"vault_key"`
	NewLimit  uint64           `borsh:"new_limit"`
	Timestamp int64            `borsh:"timestamp"`
}

// VaultUpdateDirectWithdrawEnabledEvent event struct
type VaultUpdateDirectWithdrawEnabledEvent struct {
	VaultKey                 solana.PublicKey `borsh:"vault_key"`
	NewDirectWithdrawEnabled bool             `borsh:"new_direct_withdraw_enabled"`
	Timestamp                int64            `borsh:"timestamp"`
}

// VaultUpdateMinTotalIdleEvent event struct
type VaultUpdateMinTotalIdleEvent struct {
	VaultKey        solana.PublicKey `borsh:"vault_key"`
	NewMinTotalIdle uint64           `borsh:"new_min_total_idle"`
	Timestamp       int64            `borsh:"timestamp"`
}

// VaultUpdateMinUserDepositEvent event struct
type VaultUpdateMinUserDepositEvent struct {
	VaultKey          solana.PublicKey `borsh:"vault_key"`
	NewMinUserDeposit uint64           `borsh:"new_min_user_deposit"`
	Timestamp         int64            `borsh:"timestamp"`
}

// VaultUpdateProfitMaxUnlockTimeEvent event struct
type VaultUpdateProfitMaxUnlockTimeEvent struct {
	VaultKey               solana.PublicKey `borsh:"vault_key"`
	NewProfitMaxUnlockTime uint64           `borsh:"new_profit_max_unlock_time"`
	Timestamp              int64            `borsh:"timestamp"`
}

// VaultUpdateUserDepositLimitEvent event struct
type VaultUpdateUserDepositLimitEvent struct {
	VaultKey            solana.PublicKey `borsh:"vault_key"`
	NewUserDepositLimit uint64           `borsh:"new_user_deposit_limit"`
	Timestamp           int64            `borsh:"timestamp"`
}

// VaultUpdateWhitelistedOnlyEvent event struct
type VaultUpdateWhitelistedOnlyEvent struct {
	VaultKey           solana.PublicKey `borsh:"vault_key"`
	NewWhitelistedOnly bool             `borsh:"new_whitelisted_only"`
	Timestamp          int64            `borsh:"timestamp"`
}

// VaultWithdrawlEvent event struct
type VaultWithdrawlEvent struct {
	VaultKey         solana.PublicKey `borsh:"vault_key"`
	TotalIdle        uint64           `borsh:"total_idle"`
	TotalShare       uint64           `borsh:"total_share"`
	AssetsToTransfer uint64           `borsh:"assets_to_transfer"`
	SharesToBurn     uint64           `borsh:"shares_to_burn"`
	TokenAccount     solana.PublicKey `borsh:"token_account"`
	ShareAccount     solana.PublicKey `borsh:"share_account"`
	TokenMint        solana.PublicKey `borsh:"token_mint"`
	ShareMint        solana.PublicKey `borsh:"share_mint"`
	Authority        solana.PublicKey `borsh:"authority"`
	SharePrice       uint64           `borsh:"share_price"`
	Timestamp        int64            `borsh:"timestamp"`
}

// WhitelistUpdatedEvent event struct
type WhitelistUpdatedEvent struct {
	User        solana.PublicKey `borsh:"user"`
	Whitelisted bool             `borsh:"whitelisted"`
}

// WithdrawalRequestCanceledEvent event struct
type WithdrawalRequestCanceledEvent struct {
	User      solana.PublicKey `borsh:"user"`
	Vault     solana.PublicKey `borsh:"vault"`
	Index     uint64           `borsh:"index"`
	Timestamp int64            `borsh:"timestamp"`
}

// WithdrawalRequestFulfilledEvent event struct
type WithdrawalRequestFulfilledEvent struct {
	User      solana.PublicKey `borsh:"user"`
	Vault     solana.PublicKey `borsh:"vault"`
	Amount    uint64           `borsh:"amount"`
	Index     uint64           `borsh:"index"`
	Timestamp int64            `borsh:"timestamp"`
}

// WithdrawalRequestedEvent event struct
type WithdrawalRequestedEvent struct {
	User        solana.PublicKey `borsh:"user"`
	Vault       solana.PublicKey `borsh:"vault"`
	Recipient   solana.PublicKey `borsh:"recipient"`
	Shares      uint64           `borsh:"shares"`
	Amount      uint64           `borsh:"amount"`
	MaxLoss     uint64           `borsh:"max_loss"`
	FeeShares   uint64           `borsh:"fee_shares"`
	Index       uint64           `borsh:"index"`
	Timestamp   int64            `borsh:"timestamp"`
	PriorityFee uint64           `borsh:"priority_fee"`
}
