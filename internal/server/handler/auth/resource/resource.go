package resource

import (
	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/auth"
)

func Page(c *gin.Context) (auth.Resource, error) {
	// get attributes
	// maybe need to call backend functions such as GetMetadata
	return auth.KeyValues{
		"kind": "page",
		"path": c.Request.URL.Path,
	}, nil
}

func AdminPage(c *gin.Context) (auth.Resource, error) {
	return auth.KeyValues{
		"kind": "page.admin",
		"path": c.Request.URL.Path,
	}, nil
}
func Users(c *gin.Context) (auth.Resource, error) {
	return fromParams("user", c.Params), nil
}
func Groups(c *gin.Context) (auth.Resource, error) {
	return fromParams("group", c.Params), nil
}
func Roles(c *gin.Context) (auth.Resource, error) {
	return fromParams("role", c.Params), nil
}
func Assigns(c *gin.Context) (auth.Resource, error) {
	return fromParams("assign", c.Params), nil
}
func Messages(c *gin.Context) (auth.Resource, error) {
	return auth.KeyValues{
		"kind": "message",
	}, nil
}
func fromParams(kind string, params gin.Params) auth.KeyValues {
	kvs := auth.KeyValues{
		"kind": kind,
	}
	for _, param := range params {
		kvs[param.Key] = param.Value
	}
	return kvs
}
