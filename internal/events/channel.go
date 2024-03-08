package events

type Channel struct {
	listener map[chan<- Event]struct{}
}

func (ch *Channel) fire(m Event) error {
	for l := range ch.listener {
		l <- m
	}
	return nil
}
