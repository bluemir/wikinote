package pubsub

import (
	"time"

	"github.com/sirupsen/logrus"
)

type chanEventHandler struct {
	ch chan<- Message
}

func (h chanEventHandler) Handle(ctx Context, evt Message) {
	for {
		select {
		case h.ch <- evt:
			return
		case <-time.After(10 * time.Second):
			logrus.Errorf("event handler hang detected: %#v", evt)
		}
	}
}
