package storage

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"gorm.io/datatypes"
	"gorm.io/gorm/clause"
	"time"
)

type Event struct {
	Index                int64  `gorm:"primaryKey;autoIncrement"`
	TransactionSignature string `gorm:"index;uniqueIndex:idx_txsig_logindex"`
	BlockTime            int64
	Slot                 uint64
	Name                 string
	LogIndex             int            `gorm:"uniqueIndex:idx_txsig_logindex"`
	JsonEv               datatypes.JSON `gorm:"type:jsonb"`
	Mapped               bool
	CreatedAt            time.Time `gorm:"autoCreateTime"`
}

func (g *Gorm) SaveEvent(ctx context.Context, ev Event) error {
	tx := g.DB.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&ev)

	if tx.Error != nil {
		return fmt.Errorf("failed to insert event: %w", tx.Error)
	}

	if tx.RowsAffected == 0 {
		log.Warnf("Event not inserted due to conflict (probably already exists): sig=%s logIndex=%d", ev.TransactionSignature, ev.LogIndex)
	}

	return nil
}

func (g *Gorm) MarkMapped(ctx context.Context, signature, eventName string) error {
	return g.DB.WithContext(ctx).
		Model(&Event{}).
		Where("signature = ? AND event_name = ?", signature, eventName).
		Update("mapped", true).Error
}

func (g *Gorm) LoadOrderedEvents(ctx context.Context) ([]Event, error) {
	var events []Event

	if err := g.DB.WithContext(ctx).
		Model(&Event{}).
		Order("block_time ASC, transaction_signature ASC, log_index ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch ordered events: %w", err)
	}

	return events, nil
}

func (g *Gorm) LoadEventsBySlotCursor(ctx context.Context, fromSlot uint64, slotCount int) ([]Event, error) {
	var slots []uint64

	if err := g.DB.WithContext(ctx).
		Model(&Event{}).
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

	var events []Event
	if err := g.DB.WithContext(ctx).
		Model(&Event{}).
		Where("slot IN ?", slots).
		Order("slot ASC, transaction_signature ASC, log_index ASC").
		Find(&events).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch events for slots: %w", err)
	}

	return events, nil
}
