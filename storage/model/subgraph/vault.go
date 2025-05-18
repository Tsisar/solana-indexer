package subgraph

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Vault struct {
	ID                    string      `gorm:"primaryKey;column:id"`                    // Vault address
	Token                 *Token      `gorm:"foreignKey:TokenID"`                      // Token this Vault will accrue
	TokenID               string      `gorm:"column:token_id"`                         // Token ID
	ShareToken            *Token      `gorm:"foreignKey:ShareTokenID"`                 // Token representing Shares in the Vault
	ShareTokenID          string      `gorm:"column:share_token_id"`                   // Share Token ID
	DepositLimit          string      `gorm:"column:deposit_limit;default:0"`          // The maximum amount of tokens that can be deposited in this Vault (BigInt)
	Shutdown              bool        `gorm:"column:shutdown"`                         // Is vault in shutdown
	TotalDebt             string      `gorm:"column:total_debt;default:0"`             // Total amount of assets that has been deposited in strategies (BigInt)
	TotalIdle             string      `gorm:"column:total_idle;default:0"`             // Current assets held in the vault contract (BigInt)
	MinTotalIdle          string      `gorm:"column:min_total_idle;default:0"`         // Min total idle (BigInt)
	TotalShare            string      `gorm:"column:total_share;default:0"`            // Total Share (BigInt)
	Apr                   string      `gorm:"column:apr;default:0"`                    // Annual Percentage Rate (BigDecimal → string)
	SharesSupply          string      `gorm:"column:shares_supply;default:0"`          // Current supply of Shares (BigInt)
	BalanceTokens         string      `gorm:"column:balance_tokens;default:0"`         // Balance of Tokens in the Vault and its Strategies (BigInt)
	BalanceTokensIdle     string      `gorm:"column:balance_tokens_idle;default:0"`    // Current idle Token balance (BigInt)
	Activation            string      `gorm:"column:activation;default:0"`             // Creation timestamp (BigInt)
	PerformanceFees       string      `gorm:"column:performance_fees;default:0"`       // Reported protocol fees amount for the vault (BigInt)
	TotalAllocation       string      `gorm:"column:total_allocation;default:0"`       // Total allocation after orca swap (BigDecimal)
	Accountant            *Accountant `gorm:"foreignKey:AccountantID"`                 // Accountant
	AccountantID          string      `gorm:"column:accountant_id"`                    // Accountant ID
	MinUserDeposit        string      `gorm:"column:min_user_deposit;default:0"`       // Min user deposit (BigInt)
	UserDeposit           string      `gorm:"column:user_deposit;default:0"`           // User deposit (BigInt)
	UserDepositLimit      string      `gorm:"column:user_deposit_limit;default:0"`     // User deposit limit (BigInt)
	KycVerifiedOnly       bool        `gorm:"column:kyc_verified_only"`                // KYC verified only
	DirectWithdrawEnabled bool        `gorm:"column:direct_withdraw_enabled"`          // Direct withdraw enabled
	DirectDepositEnabled  bool        `gorm:"column:direct_deposit_enabled"`           // Direct deposit enabled
	WhitelistedOnly       bool        `gorm:"column:whitelisted_only"`                 // Whitelisted only
	ProfitMaxUnlockTime   string      `gorm:"column:profit_max_unlock_time;default:0"` // Profit max unlock time (BigInt)
	CurrentSharePrice     string      `gorm:"column:current_share_price;default:0"`    // Current share price (BigInt)
	LastUpdate            string      `gorm:"column:last_update;default:0"`            // Last updated timestamp (BigInt)
	TotalPriorityFees     string      `gorm:"column:total_priority_fees;default:0"`    // Priority fees (BigInt)

	// Derived relationships
	Strategies         []*Strategy           `gorm:"foreignKey:VaultID"` // Strategies for this Vault
	Deposits           []*Deposit            `gorm:"foreignKey:VaultID"` // Token deposits into the Vault
	Withdrawals        []*Withdrawal         `gorm:"foreignKey:VaultID"` // Token withdrawals from the Vault
	HistoricalApr      []*VaultHistoricalApr `gorm:"foreignKey:VaultID"` // Historical Annual Percentage Rate
	WithdrawalRequests []*WithdrawalRequest  `gorm:"foreignKey:VaultID"` // Withdrawal requests for the Vault
}

func (*Vault) TableName() string {
	return "vaults"
}

func (v *Vault) Load(ctx context.Context, db *gorm.DB) (bool, error) {
	err := db.WithContext(ctx).
		Where("id = ?", v.ID).
		First(v).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
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
