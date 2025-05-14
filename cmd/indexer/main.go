package main

import (
	"context"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/fetcher"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/Tsisar/solana-indexer/internal/websockets"
)

var ChannelSize = 10000

func main() {
	log.Debug("Starting Solana Indexer...")
	ctx := context.Background()

	db, err := storage.InitGorm()
	if err != nil {
		log.Fatalf("Failed to init Gorm DB: %v", err)
	}
	defer db.Close()

	go fetcher.Start(ctx, db)
	go websockets.Start(ctx, db)

	select {}

}
