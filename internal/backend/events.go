package backend

import (
	"context"
	"time"

	"github.com/bluemir/wikinote/internal/backend/events"
	"github.com/bluemir/wikinote/internal/pubsub"
	"github.com/pkg/errors"
)

func (backend *Backend) GetEvents(ctx context.Context) ([]pubsub.Event, error) {
	d, err := time.ParseDuration("-24h") // TODO
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return backend.Events.List(ctx, events.Since(d))
}
