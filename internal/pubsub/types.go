package pubsub

import (
	"context"
	"encoding/gob"
	"time"
)

func init() {
	gob.Register(Message{})
}

type Message struct {
	Id     string
	At     time.Time
	Kind   string
	Detail any `gorm:"type:bytes;serializer:json"`
}

type Handler interface {
	Handle(ctx Context, msg Message)
}

type Listener chan<- Message

type Context interface {
	Context() context.Context

	Set(key any, value any)
	Get(key any) (any, bool)

	IHub
}
