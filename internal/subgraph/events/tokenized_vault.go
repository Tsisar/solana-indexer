package events

// Code generated by generate_events.go; DO NOT EDIT.

import (
	"github.com/Tsisar/solana-indexer/internal/subgraph/types"
	"github.com/gagliardetto/solana-go"
)

// StrategyReportedEvent event struct
type StrategyReportedEvent struct {
	VaultKey     solana.PublicKey
	StrategyKey  solana.PublicKey
	Gain         types.BigInt
	Loss         types.BigInt
	CurrentDebt  types.BigInt
	ProtocolFees types.BigInt
	TotalFees    types.BigInt
	TotalShares  types.BigInt
	SharePrice   types.BigInt
	Timestamp    types.BigInt
}

// UpdatedCurrentDebtForStrategyEvent event struct
type UpdatedCurrentDebtForStrategyEvent struct {
	VaultKey    solana.PublicKey
	StrategyKey solana.PublicKey
	TotalIdle   types.BigInt
	TotalDebt   types.BigInt
	NewDebt     types.BigInt
}

// VaultAddStrategyEvent event struct
type VaultAddStrategyEvent struct {
	VaultKey    solana.PublicKey
	StrategyKey solana.PublicKey
	CurrentDebt types.BigInt
	MaxDebt     types.BigInt
	LastUpdate  types.BigInt
	IsActive    bool
}

// VaultDepositEvent event struct
type VaultDepositEvent struct {
	VaultKey     solana.PublicKey
	TotalDebt    types.BigInt
	TotalIdle    types.BigInt
	TotalShare   types.BigInt
	Amount       types.BigInt
	Share        types.BigInt
	TokenAccount solana.PublicKey
	ShareAccount solana.PublicKey
	TokenMint    solana.PublicKey
	ShareMint    solana.PublicKey
	Authority    solana.PublicKey
	SharePrice   types.BigInt
	Timestamp    types.BigInt
}

// VaultEmergencyWithdrawEvent event struct
type VaultEmergencyWithdrawEvent struct {
	VaultKey            solana.PublicKey
	Recipient           solana.PublicKey
	Shares              types.BigInt
	VaultTotalShares    types.BigInt
	ScaledRatio         types.BigInt
	StrategiesProcessed types.BigInt
	Timestamp           types.BigInt
}

// VaultInitEvent event struct
type VaultInitEvent struct {
	VaultKey              solana.PublicKey
	UnderlyingToken       TokenData
	Accountant            solana.PublicKey
	ShareToken            TokenData
	DepositLimit          types.BigInt
	UserDepositLimit      types.BigInt
	MinUserDeposit        types.BigInt
	KycVerifiedOnly       bool
	DirectDepositEnabled  bool
	DirectWithdrawEnabled bool
	MinimumTotalIdle      types.BigInt
	WhitelistedOnly       bool
	ProfitMaxUnlockTime   types.BigInt
}

// VaultRemoveStrategyEvent event struct
type VaultRemoveStrategyEvent struct {
	VaultKey    solana.PublicKey
	StrategyKey solana.PublicKey
	RemovedAt   types.BigInt
}

// VaultShutDownEvent event struct
type VaultShutDownEvent struct {
	VaultKey solana.PublicKey
	Shutdown bool
}

// VaultUpdateAccountantEvent event struct
type VaultUpdateAccountantEvent struct {
	VaultKey      solana.PublicKey
	NewAccountant solana.PublicKey
	Timestamp     types.BigInt
}

// VaultUpdateDepositLimitEvent event struct
type VaultUpdateDepositLimitEvent struct {
	VaultKey  solana.PublicKey
	NewLimit  types.BigInt
	Timestamp types.BigInt
}

// VaultUpdateDirectWithdrawEnabledEvent event struct
type VaultUpdateDirectWithdrawEnabledEvent struct {
	VaultKey                 solana.PublicKey
	NewDirectWithdrawEnabled bool
	Timestamp                types.BigInt
}

// VaultUpdateMinTotalIdleEvent event struct
type VaultUpdateMinTotalIdleEvent struct {
	VaultKey        solana.PublicKey
	NewMinTotalIdle types.BigInt
	Timestamp       types.BigInt
}

// VaultUpdateMinUserDepositEvent event struct
type VaultUpdateMinUserDepositEvent struct {
	VaultKey          solana.PublicKey
	NewMinUserDeposit types.BigInt
	Timestamp         types.BigInt
}

// VaultUpdateProfitMaxUnlockTimeEvent event struct
type VaultUpdateProfitMaxUnlockTimeEvent struct {
	VaultKey               solana.PublicKey
	NewProfitMaxUnlockTime types.BigInt
	Timestamp              types.BigInt
}

// VaultUpdateUserDepositLimitEvent event struct
type VaultUpdateUserDepositLimitEvent struct {
	VaultKey            solana.PublicKey
	NewUserDepositLimit types.BigInt
	Timestamp           types.BigInt
}

// VaultUpdateWhitelistedOnlyEvent event struct
type VaultUpdateWhitelistedOnlyEvent struct {
	VaultKey           solana.PublicKey
	NewWhitelistedOnly bool
	Timestamp          types.BigInt
}

// VaultWithdrawlEvent event struct
type VaultWithdrawlEvent struct {
	VaultKey         solana.PublicKey
	TotalIdle        types.BigInt
	TotalShare       types.BigInt
	AssetsToTransfer types.BigInt
	SharesToBurn     types.BigInt
	TokenAccount     solana.PublicKey
	ShareAccount     solana.PublicKey
	TokenMint        solana.PublicKey
	ShareMint        solana.PublicKey
	Authority        solana.PublicKey
	SharePrice       types.BigInt
	Timestamp        types.BigInt
}

// WhitelistUpdatedEvent event struct
type WhitelistUpdatedEvent struct {
	User        solana.PublicKey
	Whitelisted bool
}

// WithdrawalRequestCanceledEvent event struct
type WithdrawalRequestCanceledEvent struct {
	User      solana.PublicKey
	Vault     solana.PublicKey
	Index     types.BigInt
	Timestamp types.BigInt
}

// WithdrawalRequestFulfilledEvent event struct
type WithdrawalRequestFulfilledEvent struct {
	User      solana.PublicKey
	Vault     solana.PublicKey
	Amount    types.BigInt
	Index     types.BigInt
	Timestamp types.BigInt
}

// WithdrawalRequestedEvent event struct
type WithdrawalRequestedEvent struct {
	User        solana.PublicKey
	Vault       solana.PublicKey
	Recipient   solana.PublicKey
	Shares      types.BigInt
	Amount      types.BigInt
	MaxLoss     types.BigInt
	FeeShares   types.BigInt
	Index       types.BigInt
	Timestamp   types.BigInt
	PriorityFee types.BigInt
}
