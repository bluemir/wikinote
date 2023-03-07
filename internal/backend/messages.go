package backend

import "github.com/bluemir/wikinote/internal/events"

func (backend *Backend) GetMessages(name string) ([]events.Event[Message], error) {
	return backend.hub.GetEvents(name)
}
