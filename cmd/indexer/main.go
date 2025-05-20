package main

import (
	"context"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/fetcher"
	"github.com/Tsisar/solana-indexer/core/healthchecker"
	"github.com/Tsisar/solana-indexer/core/listener"
	"github.com/Tsisar/solana-indexer/core/parser"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/subgraph"
	"time"
)

func main() {
	log.Debug("[Main] Starting Solana Indexer...")
	appCtx := context.Background()

	db, err := storage.InitGorm()
	if err != nil {
		log.Fatalf("[Main] Failed to init Gorm DB: %v", err)
	}
	defer db.Close()

	if err := storage.InitCoreModels(appCtx, db); err != nil {
		log.Fatalf("[Main] Failed to init DB: %v", err)
	}

	if err := storage.InitSubgraphModels(appCtx, db); err != nil {
		log.Fatalf("[Main] Failed to init subgraph DB: %v", err)
	}

	go func() {
		err := healthchecker.Start(appCtx, db)
		if err != nil {
			subgraph.MapError(appCtx, db, err)
			log.Fatalf("[Main] DB health check failed: %v", err)
		}
	}()

	resume := false

	for {
		ctx, cancel := context.WithCancel(context.Background())

		wsReady := make(chan struct{})
		fetchDone := make(chan struct{})
		realtimeStream := make(chan string, 1000)

		go func() {
			log.Debug("[Main] Starting WebSocket...")
			if err := listener.Start(ctx, db, wsReady, realtimeStream); err != nil {
				log.Errorf("[Main] WebSocket error: %v", err)
				cancel()
			}
		}()

		select {
		case <-ctx.Done():
			break
		case <-wsReady:
			log.Info("[Main] WS ready, starting fetcher...")
			go func() {
				if err := fetcher.Start(ctx, db, resume, fetchDone); err != nil {
					log.Errorf("[Main] Fetcher error: %v", err)
					cancel()
				}
			}()
		}

		select {
		case <-ctx.Done():
			break
		case <-fetchDone:
			log.Info("[Main] fetcher done, starting parser...")
			go func() {
				if err := parser.Start(ctx, db, resume, realtimeStream); err != nil {
					log.Errorf("[Main] Parser error: %v", err)
					cancel()
				}
			}()
		}

		<-ctx.Done()
		resume = true
		log.Info("[Main] restarting full cycle in 5 seconds...")
		time.Sleep(5 * time.Second)
	}
}
