package maping

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"gorm.io/gorm"
)

func Event(ctx context.Context, db *gorm.DB, event core.Event) error {
	pretty, _ := json.MarshalIndent(event.JsonEv, "", "  ")
	log.Debugf("Maping event index: %d, \n%s\n%s \nSlot: %d, \nSignature: %s",
		event.LogIndex, event.Name, string(pretty), event.Slot, event.TransactionSignature,
	)

	if err := updateMeta(ctx, db, event); err != nil {
		return fmt.Errorf("failed to update meta: %w", err)
	}

	if err := mapEvents(ctx, db, event); err != nil {
		return fmt.Errorf("failed to map events: %w", err)
	}
	return nil
}

func Instruction(ctx context.Context, db *gorm.DB, event core.Event) error {
	pretty, _ := json.MarshalIndent(event.JsonEv, "", "  ")
	log.Debugf("Maping event index: %d, \n%s\n%s \nSlot: %d, \nSignature: %s",
		event.LogIndex, event.Name, string(pretty), event.Slot, event.TransactionSignature,
	)

	if err := updateMeta(ctx, db, event); err != nil {
		return fmt.Errorf("failed to update meta: %w", err)
	}
	return nil
}
