package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound(c *gin.Context) {
	// TODO return not found error
	c.Error(HttpError{code: http.StatusNotFound, message: c.FullPath() + " not found"})
	c.Abort()
}

func RejectNotWritten(c *gin.Context) {
	if c.Writer.Written() {
		return
	}
	c.Error(HttpError{code: http.StatusNotFound, message: c.FullPath() + " not found"})
	c.Abort()
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
