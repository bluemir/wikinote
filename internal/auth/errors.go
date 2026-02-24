package auth

import (
	"github.com/cockroachdb/errors"
)

var (
	ErrNotImplements = errors.Errorf("not implements")
	ErrUnauthorized  = errors.Errorf("unauthorized")
	ErrForbidden     = errors.Errorf("forbidden")
)
