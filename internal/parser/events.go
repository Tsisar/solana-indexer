package parser

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/internal/events"
	"github.com/Tsisar/solana-indexer/internal/maping"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"github.com/Tsisar/solana-indexer/internal/utils"
	"github.com/gagliardetto/solana-go/rpc"
	"gorm.io/datatypes"
	"strings"
)

func parseLogs(ctx context.Context, db *storage.Gorm, sig string, tx *rpc.GetTransactionResult) error {
	timestamp := utils.BlockTime(tx.BlockTime)

	for idx, msg := range tx.Meta.LogMessages {
		if strings.HasPrefix(msg, "Program data: ") {
			if err := handleLogData(ctx, db, msg, sig, tx.Slot, timestamp, idx); err != nil {
				return err
			}
		}
	}
	return nil
}

// handleLogData decodes base64 data, then attempts each Registry decoder
// until one succeeds, using only the first-8-byte discriminator to drive matching.
func handleLogData(ctx context.Context, db *storage.Gorm, msg, sig string, slot uint64, blockTime int64, index int) error {
	// 1. Strip "Program data: " prefix and decode from base64
	rawB64 := strings.TrimPrefix(msg, "Program data: ")
	data, err := base64.StdEncoding.DecodeString(rawB64)
	if err != nil {
		return fmt.Errorf("base64 decode: %w", err)
	}

	// 2. Ensure there are at least 8 bytes for the discriminator
	if len(data) < 8 {
		log.Warnf("data too short for discriminator: %x", data)
		return nil
	}

	// 3. Extract discriminator and payload
	var disc [8]byte
	copy(disc[:], data[:8])
	payload := data[8:]

	// 4. Lookup event name by discriminator
	eventName, ok := events.Discriminators[disc]
	if !ok {
		log.Warnf("Unknown discriminator: %x", disc)
		return nil
	}

	// 5. Get decoder by event name
	decoder, ok := events.Registry[eventName]
	if !ok {
		log.Warnf("No decoder for event: %s", eventName)
		return nil
	}

	// 6. Decode the payload using the decoder
	parsed, err := decoder(payload)
	if err != nil {
		return fmt.Errorf("failed to decode %s: %w", eventName, err)
	}

	jsonVal, err := json.Marshal(parsed)
	if err != nil {
		return fmt.Errorf("failed to marshal event value: %w", err)
	}

	// 7. Store the parsed event in the database
	evRecord := storage.Event{
		TransactionSignature: sig,
		Slot:                 slot,
		LogIndex:             2000 + index,
		BlockTime:            blockTime,
		Name:                 eventName,
		JsonEv:               datatypes.JSON(jsonVal),
	}
	if err := db.SaveEvent(ctx, evRecord); err != nil {
		return fmt.Errorf("save event %s: %w", eventName, err)
	}

	// 8. Event the event
	maping.Event(evRecord)

	return nil
}
