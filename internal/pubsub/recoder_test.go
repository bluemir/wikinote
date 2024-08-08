package pubsub_test

import (
	"context"
	"sync"

	"github.com/bluemir/wikinote/internal/pubsub"
)

type Recoder struct {
	recodes []pubsub.Message
	lock    sync.RWMutex
}

func NewRecoder(ctx context.Context, hub *pubsub.Hub) *Recoder {
	recoder := Recoder{}
	ch := hub.Watch("*", ctx.Done())
	go recoder.run(ch)
	return &recoder
}
func (r *Recoder) run(ch <-chan pubsub.Message) {
	for evt := range ch {
		r.lock.Lock()
		r.recodes = append(r.recodes, evt)
		r.lock.Unlock()
	}
}

func (r *Recoder) History() []string {
	ret := []string{}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, recode := range r.recodes {
		ret = append(ret, recode.Kind)
	}
	return ret
}
