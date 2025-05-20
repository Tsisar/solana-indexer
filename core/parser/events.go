package parser

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/core/events"
	"github.com/Tsisar/solana-indexer/core/utils"
	"github.com/Tsisar/solana-indexer/storage"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"github.com/Tsisar/solana-indexer/subgraph"
	"github.com/gagliardetto/solana-go/rpc"
	"gorm.io/datatypes"
	"strings"
)

// parseLogs processes the log messages from a transaction,
// identifying any base64-encoded event logs and parsing them into structured events.
func parseLogs(ctx context.Context, db *storage.Gorm, sig string, tx *rpc.GetTransactionResult) error {
	timestamp := utils.BlockTime(tx.BlockTime)

	for idx, msg := range tx.Meta.LogMessages {
		// Look only for logs that start with the expected prefix
		if strings.HasPrefix(msg, "[parser] Program data: ") {
			if err := handleLogData(ctx, db, msg, sig, tx.Slot, timestamp, idx); err != nil {
				return err
			}
		}
	}
	log.Debugf("[parser] Parsed %d log messages for transaction %s", len(tx.Meta.LogMessages), sig)
	return nil
}

// handleLogData processes a single log message that starts with "Program data: ".
// It performs the following steps:
// 1. Removes the prefix and decodes the base64 data.
// 2. Extracts the 8-byte discriminator.
// 3. Looks up the corresponding event name and decoder function.
// 4. Decodes the event payload.
// 5. Serializes it to JSON and stores the result in the database.
func handleLogData(ctx context.Context, db *storage.Gorm, msg, sig string, slot uint64, blockTime int64, index int) error {
	// 1. Strip "Program data: " prefix and decode from base64
	rawB64 := strings.TrimPrefix(msg, "[parser] Program data: ")
	data, err := base64.StdEncoding.DecodeString(rawB64)
	if err != nil {
		return fmt.Errorf("[parser] base64 decode: %w", err)
	}

	// 2. Ensure there are at least 8 bytes for the discriminator
	if len(data) < 8 {
		log.Warnf("[parser] data too short for discriminator: %x", data)
		return nil
	}

	// 3. Extract discriminator and the remaining payload
	var disc [8]byte
	copy(disc[:], data[:8])
	payload := data[8:]

	// 4. Lookup event name by discriminator
	eventName, ok := events.Discriminators[disc]
	if !ok {
		log.Warnf("[parser] Unknown discriminator: %x", disc)
		return nil
	}

	// 5. Get the decoder function for this event
	decoder, ok := events.Registry[eventName]
	if !ok {
		log.Warnf("[parser] No decoder for event: %s", eventName)
		return nil
	}

	// 6. Decode the payload
	parsed, err := decoder(payload)
	if err != nil {
		return fmt.Errorf("[parser] failed to decode %s: %w", eventName, err)
	}

	// 7. Serialize the parsed event to JSON
	jsonVal, err := json.Marshal(parsed)
	if err != nil {
		return fmt.Errorf("[parser] failed to marshal event value: %w", err)
	}

	// 8. Save the event to the database
	evRecord := core.Event{
		TransactionSignature: sig,
		Slot:                 slot,
		LogIndex:             2000 + index, // 2000+offset to avoid collisions with other log types
		BlockTime:            blockTime,
		Name:                 eventName,
		JsonEv:               datatypes.JSON(jsonVal),
	}
	if err := db.SaveEvent(ctx, evRecord); err != nil {
		return fmt.Errorf("[parser] save event %s: %w", eventName, err)
	}

	// 9. Map event for subgraph processing
	subgraph.MapEvent(ctx, db, evRecord)

	return nil
}
