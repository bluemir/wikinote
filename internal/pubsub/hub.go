package pubsub

import (
	"context"
	"reflect"
	"time"

	"github.com/bluemir/wikinote/internal/datastruct"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
)

type IHub interface {
	Publish(ctx context.Context, detail any)
	AddHandler(kind any, handler Handler[any])
	RemoveHandler(kind any, handler Handler[any])
	Watch(kind any, done <-chan struct{}, opts ...WatchOption) <-chan Event[any]
	WatchAll(done <-chan struct{}, opts ...WatchOption) <-chan Event[any]
}

type watchConfig struct {
	bufferSize int
}

type WatchOption func(*watchConfig)

func WithBuffer(n int) WatchOption {
	return func(c *watchConfig) {
		c.bufferSize = n
	}
}

var _ IHub = (*Hub)(nil)

type Hub struct {
	handlers    datastruct.Map[reflect.Type, datastruct.Set[Handler[any]]]
	broadcaster *Broadcaster[Event[any]]
}

func NewHub(ctx context.Context) (*Hub, error) {
	return &Hub{
		broadcaster: NewBroadcaster[Event[any]](),
	}, nil
}

type keyTypeHub struct{}

var keyHub = keyTypeHub{}

func HubFrom(ctx context.Context) *Hub {
	return ctx.Value(keyHub).(*Hub)
}

func (h chanEventHandler[T]) Handle(ctx context.Context, evt Event[T]) error {
	h.ch <- evt
	return nil
}

type chanEventHandler[T any] struct {
	ch chan<- Event[T]
}

func (hub *Hub) Publish(ctx context.Context, detail any) {
	kind := reflect.TypeOf(detail)

	ctx = context.WithValue(ctx, keyHub, hub)
	evt := Event[any]{
		Context: ctx,
		Id:      xid.New().String(),
		At:      time.Now(),
		Detail:  detail,
		Kind:    kind.String(),
	}

	hub.broadcaster.Broadcast(evt)

	handlers, ok := hub.handlers.Get(kind)
	if !ok {
		return
	}

	snapshot := []Handler[any]{}
	if err := handlers.Range(func(handler Handler[any]) error {
		snapshot = append(snapshot, handler)
		return nil
	}); err != nil {
		logrus.Debug(err)
	}

	for _, handler := range snapshot {
		if err := handler.Handle(ctx, evt); err != nil {
			logrus.Debug(err)
		}
	}

}
func (hub *Hub) AddHandler(kind any, handler Handler[any]) {
	handlers, _ := hub.handlers.GetOrSet(reflect.TypeOf(kind), datastruct.NewSet[Handler[any]]())
	handlers.Add(handler)
}
func (hub *Hub) RemoveHandler(kind any, handler Handler[any]) {
	handlers, _ := hub.handlers.GetOrSet(reflect.TypeOf(kind), datastruct.NewSet[Handler[any]]())
	handlers.Remove(handler)
}
func (hub *Hub) Watch(kind any, done <-chan struct{}, opts ...WatchOption) <-chan Event[any] {
	conf := watchConfig{
		bufferSize: -1,
	}
	for _, opt := range opts {
		opt(&conf)
	}

	var ch chan Event[any]
	if conf.bufferSize >= 0 {
		ch = make(chan Event[any], conf.bufferSize)
	} else {
		ch = make(chan Event[any])
	}

	h := chanEventHandler[any]{
		ch: ch,
	}

	hub.AddHandler(kind, h)
	go func() {
		<-done
		hub.RemoveHandler(kind, h)
		close(ch)
	}()

	if conf.bufferSize >= 0 {
		return ch
	}

	// Use DynamicChan to prevent deadlock during recursive publishing
	return datastruct.DynamicChan(ch)
}
func (hub *Hub) WatchAll(done <-chan struct{}, opts ...WatchOption) <-chan Event[any] {
	return hub.broadcaster.Watch(done, opts...)
}
