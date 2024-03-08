package events

import (
	"context"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type IHub interface {
	Fire(name string, detail any) error
	GetEvents(name string) ([]Event, error)
	WatchEvents(name string, done <-chan struct{}) (<-chan Event, error)
}

var _ IHub = &Hub{}

type Event struct {
	Id     string
	At     time.Time
	Name   string
	Detail any `gorm:"type:bytes;serializer:gob"`
}

func NewHub(ctx context.Context, db *gorm.DB) (*Hub, error) {
	if err := db.AutoMigrate(Event{}); err != nil {
		return nil, err
	}

	return &Hub{
		db,
		map[string]*Channel{},
	}, nil
}

type Hub struct {
	db       *gorm.DB
	channels map[string]*Channel
}

func (hub *Hub) Fire(name string, detail any) error {
	evt := Event{
		Id:     xid.New().String(),
		At:     time.Now(),
		Name:   name,
		Detail: detail,
	}

	if err := hub.db.Save(evt).Error; err != nil {
		return err
	}

	ch, ok := hub.channels[evt.Name]
	if !ok {
		// there is no channel to broadcast event
		return nil
	}

	return ch.fire(evt)
}

func (hub *Hub) GetEvents(name string) ([]Event, error) {
	ms := []Event{}

	// TODO limit
	if err := hub.db.Where(Event{Name: name}).Find(&ms).Error; err != nil {
		return nil, err
	}
	return ms, nil
}
func (hub *Hub) WatchEvents(name string, done <-chan struct{}) (<-chan Event, error) {
	c := make(chan Event)

	if _, ok := hub.channels[name]; !ok {
		hub.channels[name] = &Channel{
			listener: map[chan<- Event]struct{}{},
		}
	}
	hub.channels[name].listener[c] = struct{}{}
	go func() {
		<-done
		delete(hub.channels[name].listener, c)
		close(c)
	}()
	return c, nil
}
