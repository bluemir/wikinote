package datastruct

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	eventQueueSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "event_queue_size",
			Help: "event queue size",
		},
	)
)

func init() {
	prometheus.MustRegister(eventQueueSize)
}

func DynamicChan[T any](in <-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)

		store := NewQueue[T]()

		for {
			eventQueueSize.Set(float64(store.Len()))

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
						eventQueueSize.Set(float64(store.Len()))

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
