package storage

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"gorm.io/gorm/clause"
)

// SaveEvent inserts or updates an event record in the database.
// If a conflict occurs on (transaction_signature, log_index), it updates all fields.
func (g *Gorm) SaveEvent(ctx context.Context, ev core.Event) error {
	tx := g.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "transaction_signature"}, {Name: "log_index"}},
			UpdateAll: true,
		}).
		Create(&ev)

	if tx.Error != nil {
		return fmt.Errorf("failed to insert or update event: %w", tx.Error)
	}

	return nil
}

// MarkMapped marks a specific event as "mapped" by setting the mapped flag to true,
// identified by its transaction signature and event name.
func (g *Gorm) MarkMapped(ctx context.Context, signature, eventName string) error {
	return g.DB.WithContext(ctx).
		Model(&core.Event{}).
		Where("transaction_signature = ? AND name = ?", signature, eventName).
		Update("mapped", true).Error
}

// LoadOrderedEvents returns all events sorted in canonical order:
// by block_time, transaction_signature, and log_index.
func (g *Gorm) LoadOrderedEvents(ctx context.Context) ([]core.Event, error) {
	var events []core.Event

	if err := g.DB.WithContext(ctx).
		Model(&core.Event{}).
		Order("block_time ASC, transaction_signature ASC, log_index ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch ordered events: %w", err)
	}

	return events, nil
}

// LoadEventsBySlotCursor loads events for the next N slots after a given starting slot.
// It first retrieves a list of slot numbers, then fetches all events belonging to those slots.
func (g *Gorm) LoadEventsBySlotCursor(ctx context.Context, fromSlot uint64, slotCount int) ([]core.Event, error) {
	var slots []uint64

	// Fetch the next N slots after the cursor
	if err := g.DB.WithContext(ctx).
		Model(&core.Event{}).
		Distinct("slot").
		Where("slot > ?", fromSlot).
		Order("slot ASC").
		Limit(slotCount).
		Pluck("slot", &slots).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch slot list: %w", err)
	}

	if len(slots) == 0 {
		return nil, nil
	}

	var events []core.Event
	if err := g.DB.WithContext(ctx).
		Model(&core.Event{}).
		Where("slot IN ?", slots).
		Order("slot ASC, transaction_signature ASC, log_index ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch events for slots: %w", err)
	}

	return events, nil
}
