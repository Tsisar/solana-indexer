package main

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/config"
	"github.com/Tsisar/solana-indexer/internal/core/fetcher"
	"github.com/Tsisar/solana-indexer/internal/core/listener"
	"github.com/Tsisar/solana-indexer/internal/core/parser"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/Tsisar/solana-indexer/internal/storage/model/core"
	"github.com/Tsisar/solana-indexer/internal/subgraph"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gorm.io/gorm"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var ready atomic.Bool
var healthy atomic.Bool

func main() {
	log.Debug("[main] Starting Solana Indexer...")

	ctx := context.Background()
	var startListener func()
	var mu sync.Mutex

	healthy.Store(true)
	ready.Store(true)

	// readiness and liveness probe server
	go readinessServer()
	// metrics server
	go metricsServer()

	g, err := storage.InitGorm()
	if err != nil {
		healthy.Store(false)
		log.Fatalf("[main] Failed to init Gorm DB: %v", err)
	}
	defer g.Close()

	// Init core models
	if err := storage.InitCoreModels(ctx, g); err != nil {
		healthy.Store(false)
		log.Fatalf("[main] Failed to init DB: %v", err)
	}

	// Init subgraph models
	if err := storage.InitSubgraphModels(g); err != nil {
		healthy.Store(false)
		log.Fatalf("[main] Failed to init subgraph DB: %v", err)
	}

	// TODO: Uncomment health checker when implemented
	//go func() {
	//	if err := healthchecker.Start(ctx, g); err != nil {
	//		subgraph.MapError(ctx, g, err)
	//		healthy.Store(false)
	//		log.Fatalf("[main] DB health check failed: %v", err)
	//	}
	//}()

	stream := make(chan string, 1000)

	go func() {
		if err := parser.Start(ctx, g, stream); err != nil {
			healthy.Store(false)
			log.Fatalf("[main] Parser error: %v", err)
		}
	}()

	fetcherDoneHandler := func() {
		log.Info("[main] Fetcher done, signaling parser to start processing historical data")
		setSynced(ctx, g.DB, true)
	}

	// Start from last saved signature if configured
	resume := config.App.ResumeFromLastSignature
	if resume {
		log.Info("[main] Resuming from last saved signature...")
	}
	// Start fetcher to fetch historical data
	if err := fetcher.Start(ctx, g, resume, fetcherDoneHandler); err != nil {
		log.Errorf("[main] Fetcher error: %v", err)
	}

	// Check if the database version matches the application version
	if err := checkVersion(ctx, g.DB); err != nil {
		log.Fatalf("[main] Version check failed: %v", err)
	}
	// Load list of signatures that are not parsed
	signatures, err := g.GetOrderedNoParsedSignatures(ctx)
	if err != nil {
		log.Errorf("[main] Failed to load signatures to parse: %v", err)
	}
	for _, signature := range signatures {
		stream <- signature
	}

	// Receive handler for WebSocket listener
	listenerReceiveHandler := func(signature string) {
		log.Infof("[main] Streaming signature: %s", signature)
		stream <- signature
	}

	// Ready handler for WebSocket listener
	listenerReadyHandler := func() {
		log.Info("[main] WebSocket stream connected, ready to receive signatures")
		subgraph.RunAggregator(ctx, g)
	}

	// Error handler for WebSocket listener
	errorHandler := func(err error) {
		msg := fmt.Sprintf("[main] WebSocket listener error: %s", err.Error())
		log.Error(msg)

		subgraph.MapError(ctx, g, err)

		if ctx.Err() != nil {
			log.Info("[main] Context canceled, not restarting listener")
			return
		}

		log.Info("[main] Attempting to reconnect WebSocket listener in 3s...")
		time.AfterFunc(3*time.Second, startListener)
	}

	startListener = func() {
		go func() {
			mu.Lock()
			defer mu.Unlock()

			log.Debug("[main] Starting WebSocket listener...")
			err := listener.Start(ctx, g, listenerReceiveHandler, listenerReadyHandler, errorHandler)
			if err != nil {
				log.Errorf("[main] listener exited with error: %v", err)

				log.Info("[main] Attempting to reconnect WebSocket listener in 30s...")
				time.AfterFunc(30*time.Second, startListener)
			}
		}()
	}
	// Start the WebSocket listener
	startListener()

	select {}
}

func readinessServer() {
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if ready.Load() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		} else {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
		}
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if healthy.Load() {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("alive"))
		} else {
			http.Error(w, "not alive", http.StatusServiceUnavailable)
		}
	})

	log.Infof("[main] Health probe server listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Errorf("[main] Probe server error: %v", err)
	}
}

func metricsServer() {
	if config.App.Metrics.Enabled {
		addr := ":" + config.App.Metrics.Port
		http.Handle("/metrics", promhttp.Handler())
		log.Infof("[main] Metrics available on %s/metrics", addr)
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Errorf("[main] Prometheus server error: %v", err)
		}
	}
}

func checkVersion(ctx context.Context, db *gorm.DB) error {
	status := core.Status{}
	if _, err := status.Load(ctx, db); err != nil {
		return fmt.Errorf("failed to load status: %w", err)
	}

	// Version updated, so we need to resync
	if status.Version != config.App.Version {
		if err := storage.TruncateEvents(db); err != nil {
			return fmt.Errorf("failed to truncate events: %w", err)
		}

		if err := storage.TruncateSubgraphTables(db); err != nil {
			return fmt.Errorf("failed to truncate subgraph tables: %w", err)
		}

		if err := storage.MarkAllTransactionsUnparsed(ctx, db); err != nil {
			return fmt.Errorf("failed to reset parsed flags: %w", err)
		}

		status.Version = config.App.Version
		if err := status.Save(ctx, db); err != nil {
			return fmt.Errorf("failed to save status: %w", err)
		}
	}
	return nil
}

func setSynced(ctx context.Context, db *gorm.DB, synced bool) {
	status := core.Status{}
	if _, err := status.Load(ctx, db); err != nil {
		log.Errorf("[main] Failed to load status: %v", err)
	}
	status.Synced = synced
	if err := status.Save(ctx, db); err != nil {
		log.Errorf("[main] Failed to set status as synced: %v", err)
	}
}
