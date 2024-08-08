package pubsub_test

import (
	"sync"

	"github.com/bluemir/wikinote/internal/pubsub"
)

type CounterHandler struct {
	lock  sync.RWMutex
	count int
}

func (h *CounterHandler) Handle(ctx pubsub.Context, evt pubsub.Message) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.count++
}
func (h *CounterHandler) GetCount() int {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.count
}
