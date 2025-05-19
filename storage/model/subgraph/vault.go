package subgraph

import (
	"context"
	"errors"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Vault struct {
	ID                    string           `gorm:"primaryKey;column:id"`                    // Vault address
	Token                 *Token           `gorm:"foreignKey:TokenID"`                      // Token this Vault will accrue
	TokenID               string           `gorm:"column:token_id"`                         // Token ID
	ShareToken            *Token           `gorm:"foreignKey:ShareTokenID"`                 // Token representing Shares in the Vault
	ShareTokenID          string           `gorm:"column:share_token_id"`                   // Share Token ID
	DepositLimit          types.BigInt     `gorm:"column:deposit_limit;default:0"`          // The maximum amount of tokens that can be deposited in this Vault (BigInt)
	Shutdown              bool             `gorm:"column:shutdown"`                         // Is vault in shutdown
	TotalDebt             types.BigInt     `gorm:"column:total_debt;default:0"`             // Total amount of assets that has been deposited in strategies (BigInt)
	TotalIdle             types.BigInt     `gorm:"column:total_idle;default:0"`             // Current assets held in the vault contract (BigInt)
	MinTotalIdle          types.BigInt     `gorm:"column:min_total_idle;default:0"`         // Min total idle (BigInt)
	TotalShare            types.BigInt     `gorm:"column:total_share;default:0"`            // Total Share (BigInt)
	Apr                   types.BigDecimal `gorm:"column:apr;default:0"`                    // Annual Percentage Rate (BigDecimal → string)
	SharesSupply          types.BigInt     `gorm:"column:shares_supply;default:0"`          // Current supply of Shares (BigInt)
	BalanceTokens         types.BigInt     `gorm:"column:balance_tokens;default:0"`         // Balance of Tokens in the Vault and its Strategies (BigInt)
	BalanceTokensIdle     types.BigInt     `gorm:"column:balance_tokens_idle;default:0"`    // Current idle Token balance (BigInt)
	Activation            types.BigInt     `gorm:"column:activation;default:0"`             // Creation timestamp (BigInt)
	PerformanceFees       types.BigInt     `gorm:"column:performance_fees;default:0"`       // Reported protocol fees amount for the vault (BigInt)
	TotalAllocation       types.BigDecimal `gorm:"column:total_allocation;default:0"`       // Total allocation after orca swap (BigDecimal)
	Accountant            *Accountant      `gorm:"foreignKey:AccountantID"`                 // Accountant
	AccountantID          string           `gorm:"column:accountant_id"`                    // Accountant ID
	MinUserDeposit        types.BigInt     `gorm:"column:min_user_deposit;default:0"`       // Min user deposit (BigInt)
	UserDeposit           types.BigInt     `gorm:"column:user_deposit;default:0"`           // User deposit (BigInt)
	UserDepositLimit      types.BigInt     `gorm:"column:user_deposit_limit;default:0"`     // User deposit limit (BigInt)
	KycVerifiedOnly       bool             `gorm:"column:kyc_verified_only"`                // KYC verified only
	DirectWithdrawEnabled bool             `gorm:"column:direct_withdraw_enabled"`          // Direct withdraw enabled
	DirectDepositEnabled  bool             `gorm:"column:direct_deposit_enabled"`           // Direct deposit enabled
	WhitelistedOnly       bool             `gorm:"column:whitelisted_only"`                 // Whitelisted only
	ProfitMaxUnlockTime   types.BigInt     `gorm:"column:profit_max_unlock_time;default:0"` // Profit max unlock time (BigInt)
	CurrentSharePrice     types.BigInt     `gorm:"column:current_share_price;default:0"`    // Current share price (BigInt)
	LastUpdate            types.BigInt     `gorm:"column:last_update;default:0"`            // Last updated timestamp (BigInt)
	TotalPriorityFees     types.BigInt     `gorm:"column:total_priority_fees;default:0"`    // Priority fees (BigInt)

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
	v.TokenID = ""
	v.ShareTokenID = ""
	v.DepositLimit = types.ZeroBigInt()
	v.Shutdown = false
	v.TotalDebt = types.ZeroBigInt()
	v.TotalIdle = types.ZeroBigInt()
	v.MinTotalIdle = types.ZeroBigInt()
	v.TotalShare = types.ZeroBigInt()
	v.Apr = types.ZeroBigDecimal()
	v.SharesSupply = types.ZeroBigInt()
	v.BalanceTokens = types.ZeroBigInt()
	v.BalanceTokensIdle = types.ZeroBigInt()
	v.Activation = types.ZeroBigInt()
	v.PerformanceFees = types.ZeroBigInt()
	v.TotalAllocation = types.ZeroBigDecimal()
	v.AccountantID = ""
	v.MinUserDeposit = types.ZeroBigInt()
	v.UserDeposit = types.ZeroBigInt()
	v.UserDepositLimit = types.ZeroBigInt()
	v.KycVerifiedOnly = false
	v.DirectWithdrawEnabled = false
	v.DirectDepositEnabled = false
	v.WhitelistedOnly = false
	v.ProfitMaxUnlockTime = types.ZeroBigInt()
	v.CurrentSharePrice = types.ZeroBigInt()
	v.LastUpdate = types.ZeroBigInt()
	v.TotalPriorityFees = types.ZeroBigInt()
}

func (v *Vault) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", v.ID).
		First(v).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		v.Init()
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (v *Vault) Save(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(v).Error
}
