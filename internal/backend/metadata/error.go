package metadata

import "errors"

var (
	ErrNotFound       = errors.New("metadata not found")
	ErrNotImplemented = errors.New("metadata not implemented")
)
