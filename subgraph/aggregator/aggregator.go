package aggregator

import (
	"context"
	"fmt"
	"github.com/Tsisar/extended-log-go/log"
	"github.com/Tsisar/solana-indexer/storage/model/subgraph"
	"github.com/Tsisar/solana-indexer/subgraph/types"
	"github.com/Tsisar/solana-indexer/subgraph/utils"
	"gorm.io/gorm"
	"time"
)

func Start(ctx context.Context, db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Hour)

	if err := aggregateAndSaveSharePrice(ctx, db, "hour"); err != nil {
		log.Errorf("[aggregator] aggregation error: %v", err)
	}
	if err := aggregateAndSaveSharePrice(ctx, db, "day"); err != nil {
		log.Errorf("[aggregator] aggregation error: %v", err)
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Debug("[aggregator] Running hourly aggregation...")
				if err := aggregateAndSaveSharePrice(ctx, db, "hour"); err != nil {
					log.Errorf("[aggregator] aggregation error: %v", err)
				}
				if err := aggregateAndSaveSharePrice(ctx, db, "day"); err != nil {
					log.Errorf("[aggregator] aggregation error: %v", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func aggregateAndSaveSharePrice(ctx context.Context, db *gorm.DB, interval string) error {
	if interval != "hour" && interval != "day" {
		return fmt.Errorf("[aggregator] unsupported interval: %s", interval)
	}

	type rawRow struct {
		VaultID    string           `gorm:"column:vault_id"`
		Timestamp  int64            `gorm:"column:timestamp_sec"`
		SharePrice types.BigDecimal `gorm:"column:share_price"`
	}

	var rows []rawRow

	query := fmt.Sprintf(`
		WITH aggregated AS (
			SELECT
				vault_id,
				EXTRACT(EPOCH FROM date_trunc('%s', to_timestamp(timestamp::numeric)))::BIGINT AS timestamp_sec,
				LAST_VALUE(share_price) OVER (
					PARTITION BY vault_id, date_trunc('%s', to_timestamp(timestamp::numeric))
					ORDER BY timestamp::numeric
					ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING
				) AS share_price
			FROM share_token_data
		)
		SELECT DISTINCT ON (vault_id, timestamp_sec)
			   vault_id,
			   timestamp_sec,
			   share_price
		FROM aggregated
		ORDER BY vault_id, timestamp_sec;
	`, interval, interval)

	if err := db.Raw(query).Scan(&rows).Error; err != nil {
		return fmt.Errorf("[aggregator] %s aggregation query failed: %w", interval, err)
	}

	for _, r := range rows {
		timestamp := *types.NewBigIntFromInt64(r.Timestamp)
		timestampStr := fmt.Sprintf("%d", r.Timestamp)
		id := utils.GenerateId(r.VaultID, timestampStr, interval)

		stat := subgraph.TokenStats{
			ID:         id,
			VaultID:    r.VaultID,
			Timestamp:  timestamp,
			SharePrice: r.SharePrice,
			Interval:   interval,
		}
		if err := stat.Save(ctx, db); err != nil {
			return fmt.Errorf("[aggregator] failed to save %s token stats: %w", interval, err)
		}
	}

	return nil
}
