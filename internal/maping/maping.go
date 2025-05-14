package maping

import (
	"context"
	"encoding/json"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/storage"
)

func Start(ctx context.Context, db *storage.Gorm, eventChannel chan []byte) {
	log.Info("Starting mapping...")
	processMessages(ctx, eventChannel)
}

func processMessages(ctx context.Context, eventChannel chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-eventChannel:
			log.Warnf("Mapping message: %s", string(message))
		}
	}
}

func Event(event storage.Event) {
	pretty, _ := json.MarshalIndent(event.JsonEv, "", "  ")

	log.Debugf("Event: %s, slot: %d, index: %d, sig: %s, ts: %d \n%s",
		event.Name, event.Slot, event.Index, event.TransactionSignature, event.BlockTime, string(pretty),
	)
}

func Instruction(event storage.Event) {
	pretty, _ := json.MarshalIndent(event.JsonEv, "", "  ")

	log.Debugf("Event: %s, slot: %d, index: %d, sig: %s, ts: %d \n%s",
		event.Name, event.Slot, event.Index, event.TransactionSignature, event.BlockTime, string(pretty),
	)
}
