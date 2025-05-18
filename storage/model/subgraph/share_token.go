package subgraph

type ShareToken struct {
	ID          string                `gorm:"primaryKey;column:id"` // Associated Token Account
	Mint        []*TokenMint          `gorm:"foreignKey:ToID"`      // Token mint account
	Burn        []*TokenBurn          `gorm:"foreignKey:FromID"`    // Token burn account
	TransferIn  []*ShareTokenTransfer `gorm:"foreignKey:ToID"`      // Token transfer account in
	TransferOut []*ShareTokenTransfer `gorm:"foreignKey:FromID"`    // Token transfer account out

	TotalMinted      string `gorm:"column:total_minted;default:0"`       // Total Minted (BigDecimal)
	TotalBurnt       string `gorm:"column:total_burnt;default:0"`        // Total Burnt (BigDecimal)
	TotalTransferIn  string `gorm:"column:total_transfer_in;default:0"`  // Total Transfer In (BigDecimal)
	TotalTransferOut string `gorm:"column:total_transfer_out;default:0"` // Total Transfer Out (BigDecimal)
	CurrentPrice     string `gorm:"column:current_price;default:0"`      // Current price of the Token (BigInt)
}

func (ShareToken) TableName() string {
	return "share_tokens"
}
