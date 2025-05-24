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

	ParserCurrentSlot = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "indexer_parser_current_slot",
			Help: "Current slot being parsed by the parser/indexer",
		},
	)
)

func init() {
	prometheus.MustRegister(FetcherCurrentSlot, ParserCurrentSlot)
}
