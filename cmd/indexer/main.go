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
	"time"
)

func main() {
	log.Debug("[Main] Starting Solana Indexer...")
	appCtx := context.Background()
	resumeFromLastSignature := config.App.ResumeFromLastSignature
	if resumeFromLastSignature {
		log.Info("Main] Resuming from last saved signature...")
	}

	gorm, err := storage.InitGorm()
	if err != nil {
		log.Fatalf("[Main] Failed to init Gorm DB: %v", err)
	}
	defer gorm.Close()

	if err := storage.InitCoreModels(appCtx, gorm, resumeFromLastSignature); err != nil {
		log.Fatalf("[Main] Failed to init DB: %v", err)
	}

	if err := storage.InitSubgraphModels(appCtx, gorm, resumeFromLastSignature); err != nil {
		log.Fatalf("[Main] Failed to init subgraph DB: %v", err)
	}

	go func() {
		if err := healthchecker.Start(appCtx, gorm); err != nil {
			subgraph.MapError(appCtx, gorm, err)
			log.Fatalf("[Main] DB health check failed: %v", err)
		}
	}()

	for {
		ctx, cancel := context.WithCancel(context.Background())
		errChan := make(chan error, 1)
		wsReady := make(chan struct{}, 1)
		fetchDone := make(chan struct{}, 1)
		parseDone := make(chan struct{}, 1)
		realtimeStream := make(chan string, 1000)

		go func() {
			log.Debug("[Main] Starting WebSocket listener...")
			if err := listener.Start(ctx, gorm, wsReady, realtimeStream, errChan); err != nil {
				errChan <- err
			}
		}()

		select {
		case err := <-errChan:
			subgraph.MapError(appCtx, gorm, err)
			log.Errorf("[Main] Listener error: %v", err)
			cancel()
			goto waitAndRestart
		case <-wsReady:
			log.Info("[Main] WS ready, starting fetcher...")
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
			log.Errorf("[Main] Fetcher error: %v", err)
			cancel()
			goto waitAndRestart
		case <-fetchDone:
			log.Info("[Main] Fetcher done, starting parser for historical data...")
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
			log.Errorf("[Main] Parser error: %v", err)
			cancel()
			goto waitAndRestart
		case <-parseDone:
			log.Info("[Main] Historical parsing complete, run aggregator, entering streaming mode")
			aggregator.Start(appCtx, gorm.DB)
			resumeFromLastSignature = true
		case <-ctx.Done():
			goto waitAndRestart
		}

		select {
		case err := <-errChan:
			subgraph.MapError(appCtx, gorm, err)
			log.Errorf("[Main] Runtime error: %v", err)
			cancel()
		case <-ctx.Done():
			cancel()
		}

	waitAndRestart:
		log.Info("[Main] Restarting full cycle in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
