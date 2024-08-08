package pubsub_test

import (
	"github.com/bluemir/wikinote/internal/pubsub"
	"github.com/sirupsen/logrus"
)

type FowardHandler struct {
	to string
}

func (h FowardHandler) Handle(ctx pubsub.Context, evt pubsub.Message) {
	logrus.Trace(evt)
	ctx.Publish(h.to, evt.Detail)
}
