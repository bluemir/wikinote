package handler

import (
	"github.com/bluemir/wikinote/internal/backend"
	"github.com/bluemir/wikinote/internal/server/middleware/auth"
)

func New(backend *backend.Backend) (*Handler, error) {
	return &Handler{
		backend: backend,
	}, nil
}

type Handler struct {
	backend *backend.Backend
}

var (
	User = auth.User
)
