package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NotFound(c *gin.Context) {
	c.Error(HttpError{code: http.StatusNotFound, message: c.FullPath() + " not found"})
	c.Abort()
}

func RejectNotWritten(c *gin.Context) {
	logger := logrus.WithField("path", c.FullPath())
	logger.Warn("called")
	if c.Writer.Written() {
		logger.Trace("skip")
		return
	}
	NotFound(c)
	logger.Trace("response not found ")
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
