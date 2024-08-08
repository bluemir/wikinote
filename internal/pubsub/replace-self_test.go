package pubsub_test

import "github.com/bluemir/wikinote/internal/pubsub"

type ReplaceSelfHandler struct{}

func (ReplaceSelfHandler) Handle(ctx pubsub.Context, evt pubsub.Message) {
	ctx.RemoveHandler("do", ReplaceSelfHandler{})
	ctx.AddHandler("do", FowardHandler{to: "done"})
}
