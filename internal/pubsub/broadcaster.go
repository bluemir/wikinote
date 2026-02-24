package pubsub

import (
	"github.com/bluemir/wikinote/internal/datastruct"
)

type Broadcaster[T any] struct {
	channels datastruct.Set[chan<- T]
	closed   chan struct{}
}

func NewBroadcaster[T any]() *Broadcaster[T] {
	return &Broadcaster[T]{
		channels: datastruct.NewSet[chan<- T](),
		closed:   make(chan struct{}),
	}
}

func (b *Broadcaster[T]) Watch(done <-chan struct{}, opts ...WatchOption) <-chan T {
	conf := watchConfig{
		bufferSize: -1,
	}
	for _, opt := range opts {
		opt(&conf)
	}

	var ch chan T
	if conf.bufferSize >= 0 {
		ch = make(chan T, conf.bufferSize)
	} else {
		ch = make(chan T)
	}

	b.channels.Add(ch)

	go func() {
		defer func() {
			b.channels.Remove(ch)
			close(ch)
		}()
		select {
		case <-done:
		case <-b.closed:
		}
	}()

	if conf.bufferSize >= 0 {
		return ch
	}

	return datastruct.DynamicChan(ch)
}

func (b *Broadcaster[T]) Broadcast(v T) {
	_ = b.channels.Range(func(ch chan<- T) error {
		select {
		case <-b.closed:
			return nil
		case ch <- v:
		}
		return nil
	})
}

func (b *Broadcaster[T]) Close() {
	close(b.closed)
}
