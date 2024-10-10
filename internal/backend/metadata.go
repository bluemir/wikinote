package backend

import (
	"context"

	"github.com/bluemir/wikinote/internal/backend/metadata"
)

func (backend *Backend) ListMetadata(ctx context.Context) ([]metadata.StoreItem, error) {
	return backend.Metadata.List(ctx)
}
