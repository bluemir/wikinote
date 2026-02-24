package events

import (
	"context"
	"time"

	"github.com/bluemir/wikinote/internal/pubsub"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IEventRecoder interface {
	List(ctx context.Context, opts ...ListOptionFn) ([]pubsub.Event[any], error)
	ListWithOption(ctx context.Context, opt *ListOption) ([]pubsub.Event[any], error)
	FindByKind(ctx context.Context, kind string, opts ...ListOptionFn) ([]pubsub.Event[any], error)
	FindByKindWithOption(ctx context.Context, kind string, opt *ListOption) ([]pubsub.Event[any], error)
}

var _ IEventRecoder = (*EventRecoder)(nil)

type EventRecoder struct {
	db *gorm.DB
}

func New(ctx context.Context, db *gorm.DB, hub *pubsub.Hub) (*EventRecoder, error) {
	if err := db.AutoMigrate(&pubsub.Event[any]{}); err != nil {
		return nil, err
	}

	go func() {
		for evt := range hub.WatchAll(ctx.Done()) {
			if err := db.Create(evt).Error; err != nil {
				logrus.Error(evt) // give up recode. just log it
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

func (m *EventRecoder) List(ctx context.Context, opts ...ListOptionFn) ([]pubsub.Event[any], error) {
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

func (m *EventRecoder) ListWithOption(ctx context.Context, opt *ListOption) ([]pubsub.Event[any], error) {
	messages := []pubsub.Event[any]{}

	if err := m.db.Limit(opt.Limit).Where("at < ? AND at > ?", opt.Before, opt.After).Find(&messages).Error; err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *EventRecoder) FindByKind(ctx context.Context, kind string, opts ...ListOptionFn) ([]pubsub.Event[any], error) {
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
func (m *EventRecoder) FindByKindWithOption(ctx context.Context, kind string, opt *ListOption) ([]pubsub.Event[any], error) {
	message := []pubsub.Event[any]{}

	if err := m.db.Limit(opt.Limit).Where("at =< ? AND at >= ? AND kind = ?", opt.Before, opt.After, kind).Find(&message).Error; err != nil {
		return nil, err
	}

	return message, nil
}
