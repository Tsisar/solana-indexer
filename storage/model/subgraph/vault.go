package subgraph

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/storage/model/generic"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
)

type Vault struct {
	ID                    string           `gorm:"primaryKey;column:id"`           // Vault address
	Token                 *Token           `gorm:"foreignKey:TokenID"`             // Token this Vault will accrue
	TokenID               string           `gorm:"column:token_id"`                // Token ID
	ShareToken            *Token           `gorm:"foreignKey:ShareTokenID"`        // Token representing Shares in the Vault
	ShareTokenID          string           `gorm:"column:share_token_id"`          // Share Token ID
	DepositLimit          types.BigInt     `gorm:"column:deposit_limit"`           // The maximum amount of tokens that can be deposited in this Vault (BigInt)
	Shutdown              bool             `gorm:"column:shutdown"`                // Is vault in shutdown
	TotalDebt             types.BigInt     `gorm:"column:total_debt"`              // Total amount of assets that has been deposited in strategies (BigInt)
	TotalIdle             types.BigInt     `gorm:"column:total_idle"`              // Current assets held in the vault contract (BigInt)
	MinTotalIdle          types.BigInt     `gorm:"column:min_total_idle"`          // Min total idle (BigInt)
	TotalShare            types.BigInt     `gorm:"column:total_share"`             // Total Share (BigInt)
	Apr                   types.BigDecimal `gorm:"column:apr"`                     // Annual Percentage Rate (BigDecimal → string)
	SharesSupply          types.BigInt     `gorm:"column:shares_supply"`           // Current supply of Shares (BigInt)
	BalanceTokens         types.BigInt     `gorm:"column:balance_tokens"`          // Balance of Tokens in the Vault and its Strategies (BigInt)
	BalanceTokensIdle     types.BigInt     `gorm:"column:balance_tokens_idle"`     // Current idle Token balance (BigInt)
	Activation            types.BigInt     `gorm:"column:activation"`              // Creation timestamp (BigInt)
	PerformanceFees       types.BigInt     `gorm:"column:performance_fees"`        // Reported protocol fees amount for the vault (BigInt)
	TotalAllocation       types.BigDecimal `gorm:"column:total_allocation"`        // Total allocation after orca swap (BigDecimal)
	Accountant            *Accountant      `gorm:"foreignKey:AccountantID"`        // Accountant
	AccountantID          string           `gorm:"column:accountant_id"`           // Accountant ID
	MinUserDeposit        types.BigInt     `gorm:"column:min_user_deposit"`        // Min user deposit (BigInt)
	UserDeposit           types.BigInt     `gorm:"column:user_deposit"`            // User deposit (BigInt)
	UserDepositLimit      types.BigInt     `gorm:"column:user_deposit_limit"`      // User deposit limit (BigInt)
	KycVerifiedOnly       bool             `gorm:"column:kyc_verified_only"`       // KYC verified only
	DirectWithdrawEnabled bool             `gorm:"column:direct_withdraw_enabled"` // Direct withdraw enabled
	DirectDepositEnabled  bool             `gorm:"column:direct_deposit_enabled"`  // Direct deposit enabled
	WhitelistedOnly       bool             `gorm:"column:whitelisted_only"`        // Whitelisted only
	ProfitMaxUnlockTime   types.BigInt     `gorm:"column:profit_max_unlock_time"`  // Profit max unlock time (BigInt)
	CurrentSharePrice     types.BigInt     `gorm:"column:current_share_price"`     // Current share price (BigInt)
	LastUpdate            types.BigInt     `gorm:"column:last_update"`             // Last updated timestamp (BigInt)
	TotalPriorityFees     types.BigInt     `gorm:"column:total_priority_fees"`     // Priority fees (BigInt)

	// Derived relationships
	Strategies         []*Strategy           `gorm:"foreignKey:VaultID"` // Strategies for this Vault
	Deposits           []*Deposit            `gorm:"foreignKey:VaultID"` // Token deposits into the Vault
	Withdrawals        []*Withdrawal         `gorm:"foreignKey:VaultID"` // Token withdrawals from the Vault
	HistoricalApr      []*VaultHistoricalApr `gorm:"foreignKey:VaultID"` // Historical Annual Percentage Rate
	WithdrawalRequests []*WithdrawalRequest  `gorm:"foreignKey:VaultID"` // Withdrawal requests for the Vault
}

func (Vault) TableName() string {
	return "vaults"
}

func (v *Vault) Init() {
	v.Token = nil
	v.TokenID = ""

	v.ShareToken = nil
	v.ShareTokenID = ""

	v.DepositLimit.Zero()
	v.Shutdown = false
	v.TotalDebt.Zero()
	v.TotalIdle.Zero()
	v.MinTotalIdle.Zero()
	v.TotalShare.Zero()
	v.Apr.Zero()
	v.SharesSupply.Zero()
	v.BalanceTokens.Zero()
	v.BalanceTokensIdle.Zero()
	v.Activation.Zero()
	v.PerformanceFees.Zero()
	v.TotalAllocation.Zero()

	v.Accountant = nil
	v.AccountantID = ""

	v.MinUserDeposit.Zero()
	v.UserDeposit.Zero()
	v.UserDepositLimit.Zero()

	v.KycVerifiedOnly = false
	v.DirectWithdrawEnabled = false
	v.DirectDepositEnabled = false
	v.WhitelistedOnly = false

	v.ProfitMaxUnlockTime.Zero()
	v.CurrentSharePrice.Zero()
	v.LastUpdate.Zero()
	v.TotalPriorityFees.Zero()

	v.Strategies = nil
	v.Deposits = nil
	v.Withdrawals = nil
	v.HistoricalApr = nil
	v.WithdrawalRequests = nil
}

func (v *Vault) GetID() string {
	return v.ID
}

func (v *Vault) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	return generic.LoadWithPreloads(ctx, db, v,
		"Strategies",
		"Deposits",
		"Withdrawals",
		"HistoricalApr",
		"WithdrawalRequests",
	)
}

func (v *Vault) Save(ctx context.Context, db *gorm.DB) error {
	return generic.Save(ctx, db, v)
}

// GetShareTokenMints returns a list of share token mints from the database.
func GetShareTokenMints(ctx context.Context, db *gorm.DB) ([]string, error) {
	var vaults []Vault

	if err := db.WithContext(ctx).Select("share_token_id").Find(&vaults).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch share token mints: %w", err)
	}

	var mints []string
	for _, vault := range vaults {
		mints = append(mints, vault.ShareTokenID)
	}

	return mints, nil
}
