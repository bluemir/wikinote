package events

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type IHub[T any] interface {
	Fire(Event[T]) error
	GetEvents(name string) ([]Event[T], error)
	WatchEvents(name string, done <-chan error) (<-chan Event[T], error)
}

type Event[T any] struct {
	Id     string
	At     time.Time
	Name   string
	Detail T `gorm:"type:bytes;serializer:gob"`
}

func (Event[T]) TableName() string {
	return "events"
}

func NewHub[T any](db *gorm.DB) (*Hub[T], error) {
	if err := db.AutoMigrate(Event[T]{}); err != nil {
		return nil, err
	}

	return &Hub[T]{
		db,
		map[string]*Channel[T]{},
	}, nil
}

type Hub[T any] struct {
	db       *gorm.DB
	channels map[string]*Channel[T]
}

func (hub *Hub[T]) Fire(evt Event[T]) error {
	evt.Id = xid.New().String()
	evt.At = time.Now()
	if err := hub.db.Save(evt).Error; err != nil {
		return err
	}

	ch, ok := hub.channels[evt.Name]
	if !ok {
		ch = &Channel[T]{
			listener: map[chan<- Event[T]]struct{}{},
		}
		hub.channels[evt.Name] = ch
	}

	return ch.fire(evt)
}

func (hub *Hub[T]) GetEvents(name string) ([]Event[T], error) {
	ms := []Event[T]{}

	// TODO limit
	if err := hub.db.Where(Event[T]{Name: name}).Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}
func (hub *Hub[T]) WatchEvents(name string, done <-chan error) (<-chan Event[T], error) {
	c := make(chan Event[T])
	hub.channels[name].listener[c] = struct{}{}
	go func() {
		<-done
		delete(hub.channels[name].listener, c)
		close(c)
	}()
	return c, nil
}
