package main

import (
	"context"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/fetcher"
	"github.com/Tsisar/solana-indexer/core/healthchecker"
	"github.com/Tsisar/solana-indexer/core/websockets"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/subgraph"
	"time"
)

func main() {
	log.Debug("Starting Solana Indexer...")
	appCtx := context.Background()

	db, err := storage.InitGorm()
	if err != nil {
		log.Fatalf("Failed to init Gorm DB: %v", err)
	}
	defer db.Close()

	if err := storage.InitCoreModels(appCtx, db); err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	if err := storage.InitSubgraphModels(appCtx, db); err != nil {
		log.Fatalf("Failed to init subgraph DB: %v", err)
	}

	go func() {
		err := healthchecker.Start(appCtx, db)
		if err != nil {
			subgraph.MapError(appCtx, db, err)
			log.Fatalf("DB health check failed: %v", err)
		}
	}()

	for {
		ctx, cancel := context.WithCancel(context.Background())

		wsReady := make(chan struct{})
		fetchDone := make(chan struct{})
		realtimeStream := make(chan string, 1000)

		//go logging(ctx, wsReady, fetchDone, realtimeStream)

		go func() {
			if err := websockets.Start(ctx, db, wsReady, realtimeStream); err != nil {
				log.Errorf("WebSocket error: %v", err)
				cancel()
			}
		}()

		select {
		case <-ctx.Done():
			break
		case <-wsReady:
			log.Info("Main: WS ready, starting fetcher...")
			go fetcher.Start(ctx, db, fetchDone)
		}

		select {
		case <-ctx.Done():
			break
		case <-fetchDone:
			log.Info("Main: fetcher done, starting parser...")
			//go runParser(ctx, realtimeStream)
		}

		<-ctx.Done()
		log.Info("Main: restarting full cycle in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	//go func() {
	//	err := websockets.Start(ctx, db)
	//	if err != nil {
	//		subgraph.MapError(ctx, db, err)
	//		log.Fatalf("Failed to start WebSocket: %v", err)
	//	}
	//}()
	//
	//fetcher.Start(ctx, db)
}

func logging(ctx context.Context, wsReady chan struct{}, fetchDone chan struct{}, realtimeStream chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-realtimeStream:
			log.Debugf("[CHAN] Received data from WebSocket: %v", message)

		case <-wsReady:
			log.Debugf("[CHAN] WebSocket is ready")

		case <-fetchDone:
			log.Debugf("[CHAN] Fetcher is done")

		}
	}
}
