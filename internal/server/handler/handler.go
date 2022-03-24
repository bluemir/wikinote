package handler

import (
	"github.com/bluemir/wikinote/internal/backend"
)

func New(backend *backend.Backend) (*Handler, error) {
	return &Handler{
		backend: backend,
	}, nil
}

type Handler struct {
	backend *backend.Backend
}
