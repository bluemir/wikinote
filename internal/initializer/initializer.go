package initializer

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Config struct {
	Key    string `gorm:"primary_key"`
	At     time.Time
	Config []byte
}

func Save(ctx context.Context, db *gorm.DB, key string, config any) error {
	gob.Register(config)

	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(config); err != nil {
		return err
	}

	if err := db.WithContext(ctx).Save(Config{
		Key:    key,
		At:     time.Now(),
		Config: buf.Bytes(),
	}).Error; err != nil {
		return err
	}
	return nil
}
func LoadOrInit(ctx context.Context, db *gorm.DB, key string, config any, initFuncs ...func(context.Context) error) error {
	gob.Register(config)
	if err := db.AutoMigrate(&Config{}); err != nil {
		return errors.WithStack(err)
	}

	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(config); err != nil {
		return err
	}

	status := Config{
		Key:    key,
		Config: buf.Bytes(),
	}
	if err := db.WithContext(ctx).FirstOrInit(&status).Error; err != nil {
		return err
	}

	if !status.At.IsZero() {
		// already has config
		buf := bytes.NewBuffer(status.Config)

		if err := gob.NewDecoder(buf).Decode(config); err != nil {
			return err
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
		return err
	}
	return nil
}
