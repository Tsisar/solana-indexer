package config

import (
	"github.com/Tsisar/extended-log-go/log"
	"github.com/joho/godotenv"
	"os"
)

var App *config

type config struct {
	EnableSignatureResume bool
	RPCEndpoint           string
	RPCWSEndpoint         string
	Programs              []string
	Postgres              postgres
	Metrics               metrics
}

type postgres struct {
	User     string
	Password string
	DB       string
	Host     string
	Port     string
}

type metrics struct {
	Enabled bool
	Port    string
}

func init() {
	if os.Getenv("RUNNING_IN_CONTAINER") != "true" {
		if err := godotenv.Load(); err == nil {
			log.Info(".env file successfully loaded")
		}
	}

	var err error
	App, err = loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
}

func loadConfig() (*config, error) {
	return &config{
		EnableSignatureResume: getBool("ENABLE_SIGNATURE_RESUME", true),
		RPCEndpoint:           getString("RPC_ENDPOINT", "https://api.mainnet-beta.solana.com"),
		RPCWSEndpoint:         getString("RPC_WS_ENDPOINT", "wss://api.mainnet-beta.solana.com"),
		Programs:              getStringSlice("PROGRAMS", []string{}),
		Postgres: postgres{
			User:     getString("POSTGRES_USER", "postgres"),
			Password: getString("POSTGRES_PASSWORD", "postgres"),
			DB:       getString("POSTGRES_DB", "indexer"),
			Host:     getString("POSTGRES_HOST", "localhost"),
			Port:     getString("POSTGRES_PORT", "5432"),
		},
		Metrics: metrics{
			Enabled: getBool("METRICS_ENABLED", false),
			Port:    getString("METRICS_PORT", "9040"),
		},
	}, nil
}
