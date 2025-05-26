package healthchecker

import (
	"context"
	"fmt"
	"github.com/Tsisar/solana-indexer/internal/storage"
	"time"
)

func Start(ctx context.Context, db *storage.Gorm) error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			err := Check(ctx, db)
			if err != nil {
				return fmt.Errorf("database health check failed: %w", err)
			}
			//log.Debugf("Database health check passed")
		}
	}
}
