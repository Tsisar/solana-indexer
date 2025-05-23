package main

import (
	"context"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/config"
	"github.com/Tsisar/solana-indexer/core/fetcher"
	"github.com/Tsisar/solana-indexer/core/healthchecker"
	"github.com/Tsisar/solana-indexer/core/listener"
	"github.com/Tsisar/solana-indexer/core/parser"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/subgraph"
	"github.com/Tsisar/solana-indexer/subgraph/aggregator"
	"net/http"
	"sync/atomic"
	"time"
)

var ready atomic.Bool
var healthy atomic.Bool

func main() {
	log.Debug("[main] Starting Solana Indexer...")
	healthy.Store(true)
	appCtx := context.Background()
	resumeFromLastSignature := config.App.ResumeFromLastSignature
	if resumeFromLastSignature {
		log.Info("Main] Resuming from last saved signature...")
	}

	gorm, err := storage.InitGorm()
	if err != nil {
		healthy.Store(false)
		log.Fatalf("[main] Failed to init Gorm DB: %v", err)
	}
	defer gorm.Close()

	// readiness and liveness probe server
	go func() {
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
	}()

	if err := storage.InitCoreModels(appCtx, gorm, resumeFromLastSignature); err != nil {
		healthy.Store(false)
		log.Fatalf("[main] Failed to init DB: %v", err)
	}

	if err := storage.InitSubgraphModels(appCtx, gorm, resumeFromLastSignature); err != nil {
		healthy.Store(false)
		log.Fatalf("[main] Failed to init subgraph DB: %v", err)
	}

	go func() {
		if err := healthchecker.Start(appCtx, gorm); err != nil {
			subgraph.MapError(appCtx, gorm, err)
			healthy.Store(false)
			log.Fatalf("[main] DB health check failed: %v", err)
		}
	}()

	for {
		ctx, cancel := context.WithCancel(context.Background())
		ready.Store(false)
		errChan := make(chan error, 1)
		wsReady := make(chan struct{}, 1)
		fetchDone := make(chan struct{}, 1)
		parseDone := make(chan struct{}, 1)
		realtimeStream := make(chan string, 1000)

		go func() {
			log.Debug("[main] Starting WebSocket listener...")
			if err := listener.Start(ctx, gorm, wsReady, realtimeStream, errChan); err != nil {
				errChan <- err
			}
		}()

		select {
		case err := <-errChan:
			subgraph.MapError(appCtx, gorm, err)
			log.Errorf("[main] Listener error: %v", err)
			cancel()
			goto waitAndRestart
		case <-wsReady:
			log.Info("[main] WS ready, starting fetcher...")
		case <-ctx.Done():
			goto waitAndRestart
		}

		go func() {
			if err := fetcher.Start(ctx, gorm, resumeFromLastSignature, fetchDone); err != nil {
				errChan <- err
			}
		}()
		select {
		case err := <-errChan:
			subgraph.MapError(appCtx, gorm, err)
			log.Errorf("[main] Fetcher error: %v", err)
			cancel()
			goto waitAndRestart
		case <-fetchDone:
			log.Info("[main] Fetcher done, starting parser for historical data...")
		case <-ctx.Done():
			goto waitAndRestart
		}

		go func() {
			if err := parser.Start(ctx, gorm, resumeFromLastSignature, parseDone, realtimeStream); err != nil {
				errChan <- err
			}
		}()
		select {
		case err := <-errChan:
			subgraph.MapError(appCtx, gorm, err)
			log.Errorf("[main] Parser error: %v", err)
			cancel()
			goto waitAndRestart
		case <-parseDone:
			log.Info("[main] Historical parsing complete, run aggregator, entering streaming mode")
			ready.Store(true)
			aggregator.Start(appCtx, gorm.DB)
			resumeFromLastSignature = true
		case <-ctx.Done():
			goto waitAndRestart
		}

		select {
		case err := <-errChan:
			subgraph.MapError(appCtx, gorm, err)
			log.Errorf("[main] Runtime error: %v", err)
			cancel()
		case <-ctx.Done():
			cancel()
		}

	waitAndRestart:
		log.Info("[main] Restarting full cycle in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
