package storage

import (
	"context"
	"github.com/Tsisar/solana-indexer/storage/model/core"
	"gorm.io/gorm/clause"
)

func (g *Gorm) SaveProgram(ctx context.Context, address string) error {
	return g.DB.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&core.Program{
			ID: address,
		}).Error
}
