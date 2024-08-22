package handler

import (
	"net/http"
	"strings"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/datastruct"
	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
)

type ListResponse[T any] struct {
	Items    []T  `json:"items"`
	Continue bool `json:"continue,omitempty"`
	Page     int  `json:"page,omitempty"`
	PerPage  int  `json:"per_page,omitempty"`
}

type KeyValues map[string]any

type RenderData struct {
	Context *gin.Context
	Data    any
}

func (rd *RenderData) User() *auth.User {
	u, _ := User(rd.Context)
	return u
}
func (rd *RenderData) IsSystemPage() bool {
	return strings.HasPrefix(rd.Context.Request.URL.Path, "/-/")
}
func (rd *RenderData) Can(verb auth.Verb, resource string) bool {
	u, _ := User(rd.Context)
	err := injector.Backends(rd.Context).Auth.Can(u, verb, datastruct.KeyValues{
		"kind": resource,
	})
	return err != nil
}

func with(c *gin.Context, data any) *RenderData {
	return &RenderData{
		Context: c,
		Data:    data,
	}
}

func HTML(path string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, path, with(c, KeyValues{}))
	}
}
