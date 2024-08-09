package initializer

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type InitializedStatus struct {
	Key string `gorm:"primary_key"`
	At  time.Time
}

func EnsureInitialize(ctx context.Context, db *gorm.DB, key string, initFuncs ...func(context.Context) error) error {
	if err := db.AutoMigrate(&InitializedStatus{}); err != nil {
		return errors.WithStack(err)
	}

	if err := db.WithContext(ctx).Take(&InitializedStatus{
		Key: key,
	}).Error; err == nil {
		return nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	// not initailized.

	for _, fn := range initFuncs {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	if err := db.WithContext(ctx).Save(InitializedStatus{
		Key: key,
		At:  time.Now(),
	}).Error; err != nil {
		return err
	}

	return nil
}
