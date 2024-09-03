package backend

import (
	"context"
	"time"

	"github.com/bluemir/wikinote/internal/backend/events"
	"github.com/bluemir/wikinote/internal/pubsub"
	"github.com/pkg/errors"
)

func (backend *Backend) GetMessages(ctx context.Context) ([]pubsub.Message, error) {
	d, err := time.ParseDuration("-24h") // TODO
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return backend.events.List(ctx, events.Since(d))
}
