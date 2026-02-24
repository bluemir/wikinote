package errs

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/backend/metadata"
)

func code(err error) int {
	logrus.Tracef("%T", err)

	var validationErrors validator.ValidationErrors
	var sqliteErr sqlite3.Error
	var httpErr HttpError

	// errors.Is checks same value, but errors.As checks only its type.
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return http.StatusConflict
	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound
	case errors.As(err, &validationErrors):
		return http.StatusBadRequest
	case errors.Is(err, auth.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, auth.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, os.ErrNotExist):
		return http.StatusNotFound
	case errors.Is(err, metadata.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, metadata.ErrNotImplemented):
		return http.StatusNotImplemented
	case errors.As(err, &sqliteErr):
		switch {
		case sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique:
			return http.StatusConflict
		default:
			return http.StatusNotImplemented
		}
	case errors.As(err, &httpErr):
		return httpErr.Code()
	}

	// finally check string match
	logrus.Trace(err.Error())
	switch {
	case strings.HasPrefix(err.Error(), "html/template: ") && strings.HasSuffix(err.Error(), " is undefined"):
		//html/template: ".*" is undefined
		return http.StatusNotImplemented
	}

	return http.StatusInternalServerError
}
