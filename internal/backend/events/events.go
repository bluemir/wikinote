package events

import (
	"context"
	"time"

	"github.com/bluemir/wikinote/internal/pubsub"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEventRecoder interface {
	List(ctx context.Context, opts ...ListOptionFn) ([]pubsub.Message, error)
	ListWithOption(ctx context.Context, opt *ListOption) ([]pubsub.Message, error)
	FindByKind(ctx context.Context, kind string, opts ...ListOptionFn) ([]pubsub.Message, error)
	FindByKindWithOption(ctx context.Context, kind string, opt *ListOption) ([]pubsub.Message, error)
}

var _ IEventRecoder = (*EventRecoder)(nil)

type EventRecoder struct {
	db *gorm.DB
}

func New(ctx context.Context, db *gorm.DB, hub *pubsub.Hub) (*EventRecoder, error) {
	if err := db.AutoMigrate(&pubsub.Message{}); err != nil {
		return nil, err
	}

	go func() {
		for msg := range hub.Watch("*", ctx.Done()) {
			if msg.Kind == "error" {
				logrus.Error(msg.Detail) // give up recode. just log it
				return
			}
			if err := db.Create(msg).Error; err != nil {
				hub.Publish("error", err)
			}
		}
	}()
	return &EventRecoder{
		db: db,
	}, nil
}

type ListOption struct {
	Limit  int
	After  time.Time
	Before time.Time
}
type ListOptionFn func(opt *ListOption)

func Limit(n int) ListOptionFn {
	return func(opt *ListOption) {
		opt.Limit = n
	}
}
func Since(d time.Duration) ListOptionFn {
	return func(opt *ListOption) {
		opt.After = time.Now().Add(d)
	}
}
func Until(d time.Duration) ListOptionFn {
	return func(opt *ListOption) {
		opt.Before = time.Now().Add(d)
	}
}

func (m *EventRecoder) List(ctx context.Context, opts ...ListOptionFn) ([]pubsub.Message, error) {
	opt := ListOption{
		Limit:  -1,
		After:  time.Time{},
		Before: time.Now(),
	}

	for _, fn := range opts {
		fn(&opt)
	}

	return m.ListWithOption(ctx, &opt)
}

func (m *EventRecoder) ListWithOption(ctx context.Context, opt *ListOption) ([]pubsub.Message, error) {
	message := []pubsub.Message{}

	if err := m.db.Limit(opt.Limit).Where("at < ? AND at > ?", opt.Before, opt.After).Find(&message).Error; err != nil {
		return nil, err
	}

	return message, nil
}

func (m *EventRecoder) FindByKind(ctx context.Context, kind string, opts ...ListOptionFn) ([]pubsub.Message, error) {
	opt := ListOption{
		Limit:  -1,
		After:  time.Time{},
		Before: time.Now(),
	}

	for _, fn := range opts {
		fn(&opt)
	}

	return m.ListWithOption(ctx, &opt)
}
func (m *EventRecoder) FindByKindWithOption(ctx context.Context, kind string, opt *ListOption) ([]pubsub.Message, error) {
	message := []pubsub.Message{}

	if err := m.db.Limit(opt.Limit).Where("at =< ? AND at >= ? AND kind = ?", opt.Before, opt.After, kind).Find(&message).Error; err != nil {
		return nil, err
	}

	return message, nil
}
