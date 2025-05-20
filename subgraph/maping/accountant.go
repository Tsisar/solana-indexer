package maping

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/Tsisar/solana-indexer/subgraph/events"
	"github.com/Tsisar/solana-indexer/subgraph/library/accountant"
	"gorm.io/gorm"
)

func mapEntryFeeUpdatedEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] EntryFeeUpdatedEvent: %s", event.TransactionSignature)
	var ev events.EntryFeeUpdatedEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode EntryFeeUpdatedEvent: %w", err)
	}
	if err := accountant.SetEntryFee(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to set entry fee: %w", err)
	}
	return nil
}

func mapPerformanceFeeUpdatedEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] PerformanceFeeUpdatedEvent: %s", event.TransactionSignature)
	var ev events.PerformanceFeeUpdatedEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode PerformanceFeeUpdatedEvent: %w", err)
	}
	if err := accountant.SetPerformanceFee(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to set performance fee: %w", err)
	}
	return nil
}

func mapRedemptionFeeUpdatedEvent(ctx context.Context, db *gorm.DB, event core.Event) error {
	log.Infof("[mapping] RedemptionFeeUpdatedEvent: %s", event.TransactionSignature)
	var ev events.RedemptionFeeUpdatedEvent
	if err := json.Unmarshal(event.JsonEv, &ev); err != nil {
		return fmt.Errorf("failed to decode RedemptionFeeUpdatedEvent: %w", err)
	}
	if err := accountant.SetRedemptionFee(ctx, db, ev); err != nil {
		return fmt.Errorf("failed to set entry fee: %w", err)
	}
	return nil
}
