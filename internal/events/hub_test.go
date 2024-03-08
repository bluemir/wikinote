package events

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewTestHub(ctx context.Context) (IHub, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		return nil, err
	}
	return NewHub(ctx, db)
}

func TestBusCall(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	hub, err := NewTestHub(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	ch, err := hub.WatchEvents("other", ctx.Done())
	if err != nil {
		t.Fatal(err)
	}

	counter := runEventHandler(ch)

	hub.Fire("other", struct{}{})
	hub.Fire("other", struct{}{})
	hub.Fire("other", struct{}{})

	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, 3, *counter)
}

func runEventHandler(ch <-chan Event) *int {
	c := 0
	go func() {
		for evt := range ch {
			c += 1

			logrus.Trace(evt)
		}
	}()
	return &c
}
