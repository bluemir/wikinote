package pubsub

import (
	"context"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/datastruct"
)

type IRouter[T any] interface {
	Publish(ctx context.Context, kind string, detail T)
	AddHandler(kind string, handler Handler[T])
	RemoveHandler(kind string, handler Handler[T])
	Watch(kind string, done <-chan struct{}, opts ...WatchOption) <-chan Event[T]
	WatchAll(done <-chan struct{}, opts ...WatchOption) <-chan Event[T]
}

var _ IRouter[any] = (*Router[any])(nil)

type Router[T any] struct {
	handlers    datastruct.Tree[string, datastruct.Set[Handler[T]]]
	broadcaster *Broadcaster[Event[T]]
}

const Separator = "." // QUESTION make configurable?

func NewRoute[T any](ctx context.Context) (*Router[T], error) {
	return &Router[T]{
		broadcaster: NewBroadcaster[Event[T]](),
	}, nil
}

type keyTypeRouter struct{}

var keyRouter = keyTypeRouter{}

func RouterFrom[T any](ctx context.Context) *Router[T] {
	return ctx.Value(keyRouter).(*Router[T])
}

func (router *Router[T]) Publish(ctx context.Context, kind string, detail T) {
	keys := strings.Split(kind, Separator)
	handlers, ok := router.handlers.Get(keys...)

	ctx = context.WithValue(ctx, keyRouter, router)
	evt := Event[T]{
		Context: ctx,
		Id:      xid.New().String(),
		At:      time.Now(),
		Detail:  detail,
		Kind:    kind,
	}

	if ok {
		snapshot := []Handler[T]{}
		if err := handlers.Range(func(handler Handler[T]) error {
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

	router.broadcaster.Broadcast(evt)
}
func (router *Router[T]) AddHandler(kind string, handler Handler[T]) {
	keys := strings.Split(kind, Separator)
	handlers, _ := router.handlers.GetOrSet(keys, datastruct.NewSet[Handler[T]]())

	handlers.Add(handler)
}
func (router *Router[T]) RemoveHandler(kind string, handler Handler[T]) {
	keys := strings.Split(kind, Separator)
	handlers, _ := router.handlers.GetOrSet(keys, datastruct.NewSet[Handler[T]]())

	handlers.Remove(handler)
}
func (router *Router[T]) Watch(kind string, done <-chan struct{}, opts ...WatchOption) <-chan Event[T] {
	conf := watchConfig{
		bufferSize: -1,
	}
	for _, opt := range opts {
		opt(&conf)
	}

	var ch chan Event[T]
	if conf.bufferSize >= 0 {
		ch = make(chan Event[T], conf.bufferSize)
	} else {
		ch = make(chan Event[T])
	}

	h := chanEventHandler[T]{
		ch: ch,
	}

	router.AddHandler(kind, h)
	go func() {
		<-done
		router.RemoveHandler(kind, h)
		close(ch)
	}()

	if conf.bufferSize >= 0 {
		return ch
	}

	// Use DynamicChan to prevent deadlock during recursive publishing
	return datastruct.DynamicChan(ch)
}
func (router *Router[T]) WatchAll(done <-chan struct{}, opts ...WatchOption) <-chan Event[T] {
	return router.broadcaster.Watch(done, opts...)
}
