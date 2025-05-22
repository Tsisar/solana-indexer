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

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Debug("Running hourly aggregation...")
				if err := AggregateAndSaveHourlySharePrice(ctx, db); err != nil {
					log.Errorf("Aggregation error: %v", err)
				}
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func AggregateAndSaveHourlySharePrice(ctx context.Context, db *gorm.DB) error {
	type rawRow struct {
		VaultID    string           `gorm:"column:vault_id"`
		Timestamp  int64            `gorm:"column:timestamp_sec"`
		SharePrice types.BigDecimal `gorm:"column:share_price"`
	}

	var rows []rawRow

	query := `
		WITH aggregated AS (
			SELECT
				vault_id,
				EXTRACT(EPOCH FROM date_trunc('day', to_timestamp(timestamp::numeric)))::BIGINT AS timestamp_sec,
				LAST_VALUE(share_price) OVER (
					PARTITION BY vault_id, date_trunc('day', to_timestamp(timestamp::numeric))
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
	`

	if err := db.Raw(query).Scan(&rows).Error; err != nil {
		return fmt.Errorf("daily aggregation query failed: %w", err)
	}

	for _, r := range rows {
		timestamp := *types.NewBigIntFromInt64(r.Timestamp)
		timestampStr := fmt.Sprintf("%d", r.Timestamp)
		id := utils.GenerateId(r.VaultID, timestampStr)

		stat := subgraph.TokenStats{
			ID:         id,
			VaultID:    r.VaultID,
			Timestamp:  timestamp,
			SharePrice: r.SharePrice,
		}
		if err := stat.Save(ctx, db); err != nil {
			return fmt.Errorf("failed to save token stats: %w", err)
		}
	}

	return nil
}
