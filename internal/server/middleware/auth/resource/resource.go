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
func Global(c *gin.Context) (auth.Resource, error) {
	return auth.KeyValues{}, nil
}
func AdminPage(c *gin.Context) (auth.Resource, error) {
	return auth.KeyValues{
		"kind": "admin",
	}, nil
}
