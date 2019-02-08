package plugins

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/pkgs/renderer"
)

type RenderClause interface {
	With(*gin.Context)
}

func RenderPage(data map[string]interface{}) RenderClause {
	return Meta(data)
}

type Meta map[string]interface{}

func (rc Meta) With(c *gin.Context) {
	c.HTML(http.StatusOK, "/plugin-template.html", renderer.Data(rc).With(c))
}
