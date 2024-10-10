package pubsub

import (
	"context"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/datastruct"
)

func NewHub(ctx context.Context) (*Hub, error) {
	in := make(chan Message)
	go func() {
		<-ctx.Done()
		close(in)
	}()

	q := datastruct.DynamicChan(in)

	hub := &Hub{
		ctx:      ctx,
		values:   datastruct.Map[any, any]{},
		handlers: datastruct.NewTree[string, datastruct.Set[Handler]](),
		in:       in,
	}

	go func() {
		for evt := range q {
			hub.broadcast(evt)
		}
	}()

	return hub, nil
}

type Hub struct {
	ctx      context.Context
	values   datastruct.Map[any, any]
	handlers *datastruct.Tree[string, datastruct.Set[Handler]]

	in chan<- Message
}

func (h *Hub) Publish(kind string, detail any) {
	logrus.Tracef("fire event: %s - %#v", kind, detail)
	if strings.Contains(kind, "*") {
		return // error or send error as messagez?
	}

	h.in <- Message{
		Id:     xid.New().String(),
		At:     time.Now(),
		Kind:   kind,
		Detail: detail,
	}
}
func (h *Hub) broadcast(evt Message) {
	logrus.Tracef("broadcast event: %s - %#v", evt.Kind, evt.Detail)

	keys := strings.Split(evt.Kind, ".")
	{
		handlers, _ := h.handlers.GetOrSet(keys, datastruct.NewSet[Handler]())

		logrus.Tracef("handlers: %d", handlers.Len())

		for _, handler := range handlers.List() {
			handler.Handle(h, evt)
		}
	}
	{
		// handle star
		handlers, _ := h.handlers.GetOrSet([]string{"*"}, datastruct.NewSet[Handler]())
		for _, handler := range handlers.List() {
			handler.Handle(h, evt)
		}
	}
}
func (h *Hub) Close() {
	close(h.in)
}

func (h *Hub) AddHandler(kind string, handler Handler) {
	if handler == nil {
		return
	}
	keys := strings.Split(kind, ".")
	set, _ := h.handlers.GetOrSet(keys, datastruct.NewSet[Handler]())

	set.Add(handler)
}
func (h *Hub) RemoveHandler(kind string, handler Handler) {
	keys := strings.Split(kind, ".")
	set, _ := h.handlers.GetOrSet(keys, datastruct.NewSet[Handler]())

	set.Remove(handler)
}

func (h *Hub) AddListener(kind string, l Listener) {
	logrus.Tracef("register listener: %s", kind)
	h.AddHandler(kind, chanEventHandler{l})
}
func (h *Hub) RemoveListener(kind string, l Listener) {
	h.RemoveHandler(kind, chanEventHandler{l})
}

func (h *Hub) Watch(kind string, done <-chan struct{}) <-chan Message {
	logrus.WithField("kind", kind).Trace("watch event")

	c := make(chan Message)
	h.AddListener(kind, c)
	go func() {
		<-done
		h.RemoveListener(kind, c)
		close(c)
		logrus.WithField("kind", kind).Trace("unwatch event")
	}()

	return c
}

func (h *Hub) Context() context.Context {
	return h.ctx
}

func (h *Hub) Set(key any, value any) {
	h.values.Set(key, value)
}
func (h *Hub) Get(key any) (any, bool) {
	return h.values.Get(key)
}
