package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}
func (e NotFoundError) Code() int {
	return http.StatusNotFound
}
func NotFound(c *gin.Context) {
	// TODO return not found error
	c.Error(NotFoundError(c.FullPath() + " not found"))
	c.Abort()
}

type ForbiddenError string

func (e ForbiddenError) Error() string {
	return string(e)
}
func (e ForbiddenError) Code() int {
	return http.StatusForbidden
}
