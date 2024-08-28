package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func NotFound(c *gin.Context) {
	// TODO return not found error
	c.Error(HttpError{code: http.StatusNotFound, message: c.FullPath() + " not found"})
	c.Abort()
}

func NotFoundWithPrefix(prefixs ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				c.Error(HttpError{code: http.StatusNotFound, message: c.FullPath() + " not found"})
				c.Abort()
				return
			}
		}
	}
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
