package pubsub

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

type chanEventHandler struct {
	ch chan<- Event
}

func (h chanEventHandler) Handle(ctx context.Context, evt Event) {
	for {
		select {
		case h.ch <- evt:
			return
		case <-time.After(10 * time.Second):
			logrus.Errorf("event handler hang detected: %#v", evt)
		}
	}
}
