package events

type Channel[T any] struct {
	listener map[chan<- Event[T]]struct{}
}

func (ch *Channel[T]) fire(m Event[T]) error {
	for l := range ch.listener {
		l <- m
	}
	return nil
}
