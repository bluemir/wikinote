package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
)

func ListMetadata(c *gin.Context) error {
	items, err := injector.Backends(c).ListMetadata(c.Request.Context())
	if err != nil {
		return err
	}

	c.HTML(http.StatusOK, "admin/metadata.html", With(c, gin.H{
		"Items": items,
	}))
	return nil
}
