package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
)

func ListMetadata(c *gin.Context) {
	items, err := injector.Backends(c).ListMetadata(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "admin/metadata.html", With(c, gin.H{
		"Items": items,
	}))
}
