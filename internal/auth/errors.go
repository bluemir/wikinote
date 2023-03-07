package auth

import (
	"github.com/pkg/errors"
)

var (
	ErrNotImplements = errors.Errorf("not implements")
	ErrUnauthorized  = errors.Errorf("unauthorized")
	ErrForbidden     = errors.Errorf("forbidden")
)
