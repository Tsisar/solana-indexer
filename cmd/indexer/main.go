package main

import (
	"context"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/fetcher"
	"github.com/Tsisar/solana-indexer/core/healthchecker"
	"github.com/Tsisar/solana-indexer/core/websockets"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/subgraph"
)

func main() {
	log.Debug("Starting Solana Indexer...")
	ctx := context.Background()

	db, err := storage.InitGorm()
	if err != nil {
		log.Fatalf("Failed to init Gorm DB: %v", err)
	}
	defer db.Close()

	if err := storage.InitCoreModels(ctx, db); err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	if err := storage.InitSubgraphModels(ctx, db); err != nil {
		log.Fatalf("Failed to init subgraph DB: %v", err)
	}

	go func() {
		err := healthchecker.Start(ctx, db)
		if err != nil {
			subgraph.MapError(ctx, db, err)
			log.Fatalf("DB health check failed: %v", err)
		}
	}()

	go func() {
		err := websockets.Start(ctx, db)
		if err != nil {
			subgraph.MapError(ctx, db, err)
			log.Fatalf("Failed to start WebSocket: %v", err)
		}
	}()

	fetcher.Start(ctx, db)

	select {}
}
