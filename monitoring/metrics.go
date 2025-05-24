package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
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
	)
}
