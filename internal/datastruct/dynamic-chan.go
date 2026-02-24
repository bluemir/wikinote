package datastruct

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	EventQueueSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "event_queue_size",
			Help: "event queue size",
		},
	)

	// if want to know queue size, run this.
	// prometheus.MustRegister(EventQueueSize)
)

func DynamicChan[T any](in <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		store := NewQueue[T]()

		for {
			EventQueueSize.Set(float64(store.Len()))

			if store.Len() == 0 {
				evt, more := <-in
				if !more {
					return
				}
				store.Add(evt)
				continue
			}

			select {
			case evt, more := <-in:
				if !more {
					for store.Len() > 0 {
						EventQueueSize.Set(float64(store.Len()))

						out <- store.Front()
						store.Pop()
					}
					return
				}
				store.Add(evt)
			case out <- store.Front():
				store.Pop()
			}
		}
	}()

	return out
}
