package monitoring

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"math/big"
)

var (
	FetcherCurrentSlot = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "indexer_fetcher_current_slot",
			Help: "Current slot being fetched by the fetcher",
		},
	)

	ListenerCurrentSlot = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "indexer_listener_current_slot",
			Help: "Current slot received by the WebSocket listener",
		},
	)

	ParserCurrentSlot = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "indexer_parser_current_slot",
			Help: "Current slot being parsed by the parser/indexer",
		},
	)

	DepositsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "indexer_deposit_total",
			Help: "Number of processed deposits",
		},
		[]string{"vault_id", "token_id"},
	)

	WithdrawalsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "indexer_withdrawal_total",
			Help: "Number of processed withdrawals",
		},
		[]string{"vault_id", "token_id"},
	)

	DepositTokenSum = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "indexer_deposit_token_sum",
			Help: "Total token amount deposited to each vault",
		},
		[]string{"vault_id", "token_id"},
	)

	WithdrawalTokenSum = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "indexer_withdrawal_token_sum",
			Help: "Total token amount withdrawn from each vault",
		},
		[]string{"vault_id", "token_id"},
	)

	TokenPrice = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "indexer_token_price",
			Help: "Current token price",
		},
		[]string{"token_id", "symbol", "name"},
	)

	TokenDecimals = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "indexer_token_decimals",
			Help: "Token decimals",
		},
		[]string{"token_id"},
	)
)

func init() {
	prometheus.MustRegister(
		FetcherCurrentSlot,
		ParserCurrentSlot,
		ListenerCurrentSlot,
		DepositsTotal,
		WithdrawalsTotal,
		DepositTokenSum,
		WithdrawalTokenSum,
		TokenPrice,
		TokenDecimals,
	)
}

func Withdrawal(ctx context.Context, db *gorm.DB, withdraw subgraph.Withdrawal) {
	token := subgraph.Token{ID: withdraw.TokenID}
	ok, err := token.Load(ctx, db)
	if err != nil || !ok {
		return
	}

	// amount / 10^decimals using big.Float
	amountFloat := new(big.Float).SetInt(withdraw.TokenAmount.Int)
	tenPow := new(big.Int).Exp(big.NewInt(10), token.Decimals.Int, nil)
	normalized := new(big.Float).Quo(amountFloat, new(big.Float).SetInt(tenPow))

	f64, _ := normalized.Float64()

	WithdrawalsTotal.WithLabelValues(withdraw.VaultID, withdraw.TokenID).Inc()
	WithdrawalTokenSum.WithLabelValues(withdraw.VaultID, withdraw.TokenID).Add(f64)
}

func Deposit(ctx context.Context, db *gorm.DB, deposit subgraph.Deposit) {
	token := subgraph.Token{ID: deposit.TokenID}
	ok, err := token.Load(ctx, db)
	if err != nil || !ok {
		return
	}

	// amount / 10^decimals using big.Float
	amountFloat := new(big.Float).SetInt(deposit.TokenAmount.Int)
	tenPow := new(big.Int).Exp(big.NewInt(10), token.Decimals.Int, nil)
	normalized := new(big.Float).Quo(amountFloat, new(big.Float).SetInt(tenPow))

	f64, _ := normalized.Float64()

	DepositsTotal.WithLabelValues(deposit.VaultID, deposit.TokenID).Inc()
	DepositTokenSum.WithLabelValues(deposit.VaultID, deposit.TokenID).Add(f64)
}

func Token(t subgraph.Token) {
	priceFloat, _ := t.CurrentPrice.Float64()
	decimals, _ := t.Decimals.Float64()

	TokenPrice.WithLabelValues(t.ID, t.Symbol, t.Name).Set(priceFloat)
	TokenDecimals.WithLabelValues(t.ID).Set(decimals)
}
