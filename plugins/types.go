package plugins

import (
	"html/template"

	"github.com/gin-gonic/gin"
)

type FooterPlugin interface {
	Footer(path string) (template.HTML, error)
}
type AfterWikiSavePlugin interface {
	AfterWikiSave(path string, data []byte) error
}
type RegisterRouterPlugin interface {
	RegisterRouter(r gin.IRouter)
}
