package pubsub

import (
	"context"
	"time"
)

type Event struct {
	Context context.Context `gorm:"-" json:"-"`
	Id      string
	At      time.Time
	Detail  any `gorm:"type:bytes;serializer:gob"`
	Kind    string
	// Event 에 Kind 를 넣어야 할까?
	// 보통은 kind 를 지정 해서 watch 하는 것은 detail의 type 이 결정 되는 것인데, watch all 의 경우만 그렇지 않다.
}

type Handler interface {
	Handle(ctx context.Context, evt Event)
}
