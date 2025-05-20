package generic

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type Identifiable interface {
	GetID() string
	Init()
}

func Load[T Identifiable](ctx context.Context, db *gorm.DB, model T) (bool, error) {
	err := db.WithContext(ctx).
		Session(&gorm.Session{Logger: db.Logger.LogMode(logger.Silent)}).
		Where("id = ?", model.GetID()).
		First(model).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		model.Init()
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func Save[T any](ctx context.Context, db *gorm.DB, model T) error {
	return db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(model).Error
}
