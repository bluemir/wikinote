package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/gin-gonic/gin"
)

type ListResponse[T any] struct {
	Items    []T  `json:"items"`
	Continue bool `json:"continue,omitempty"`
	Page     int  `json:"page,omitempty"`
	PerPage  int  `json:"per_page,omitempty"`
}

type KeyValues map[string]any

func renderData(c *gin.Context, data any) KeyValues {
	return KeyValues{
		"context": c,
		"data":    data,
		"user": func() *auth.User {
			u, _ := User(c)
			return u
		},
	}
}

func HTML(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, path, renderData(c, KeyValues{}))
	}
}
