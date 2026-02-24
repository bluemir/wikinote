package handler

import (
	"github.com/cockroachdb/errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	ErrNotFile          = errors.New("target is not a file")
	ErrNotDirectory     = errors.New("target is not a directory")
	ErrPasswordMismatch = errors.New("password and confirm password not matched")
	ErrInvalidCode      = errors.New("invalid code")
)

func NotFound(c *gin.Context) error {
	return HttpError{code: http.StatusNotFound, message: c.FullPath() + " not found"}
}

func RejectNotWritten(c *gin.Context) error {
	logger := logrus.WithField("path", c.FullPath())
	logger.Warn("called")
	if c.Writer.Written() {
		logger.Trace("skip")
		return nil
	}
	logger.Trace("response not found ")
	return NotFound(c)
}

type HttpError struct {
	code    int
	message string
}

func (e HttpError) Error() string {
	return e.message
}
func (e HttpError) Code() int {
	return e.code
}
