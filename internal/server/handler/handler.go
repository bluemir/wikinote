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

type ListResponse[T any] struct {
	Items    []T  `json:"items"`
	Continue bool `json:"continue,omitempty"`
	Page     int  `json:"page,omitempty"`
	PerPage  int  `json:"per_page,omitempty"`
}
