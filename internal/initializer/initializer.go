package initializer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"
	"gorm.io/gorm"
)

type Config struct {
	Key    string `gorm:"primary_key"`
	At     time.Time
	Config []byte
}

func Save(ctx context.Context, db *gorm.DB, key string, config any) error {
	buf, err := json.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := db.WithContext(ctx).Save(Config{
		Key:    key,
		At:     time.Now(),
		Config: buf,
	}).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
func LoadOrInit(ctx context.Context, db *gorm.DB, key string, config any, initFuncs ...func(context.Context) error) error {
	if err := db.AutoMigrate(&Config{}); err != nil {
		return errors.WithStack(err)
	}

	buf, err := json.Marshal(config)
	if err != nil {
		return errors.WithStack(err)
	}

	status := Config{
		Key:    key,
		Config: buf,
	}
	if err := db.WithContext(ctx).Where(Config{
		Key: key,
	}).FirstOrInit(&status).Error; err != nil {
		return errors.WithStack(err)
	}

	if !status.At.IsZero() {
		// already has config
		if err := json.Unmarshal(status.Config, config); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}

	// not initialized
	status.At = time.Now()
	for _, fn := range initFuncs {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	if err := db.WithContext(ctx).Create(&status).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
