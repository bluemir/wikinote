package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/server/injector"
)

func ListPlugins(c *gin.Context) {
	plugins, err := injector.Backends(c).Plugin.List(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "admin/plugins/list.html", With(c, plugins))
}
func GetPlugin(c *gin.Context) {
	plugin, err := injector.Backends(c).Plugin.Get(c.Request.Context(), c.Param("name"))
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "admin/plugins/config.html", With(c, plugin))
}
func UpdatePlugin(c *gin.Context) {
	req := struct {
		IsEnabled bool   `form:"isEnabled"`
		Config    string `form:"config"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		return
	}

	name := c.Param("name")

	if req.IsEnabled {
		if err := injector.Backends(c).Plugin.Enable(c.Request.Context(), name); err != nil {
			c.Error(err)
			return
		}
	} else {
		if err := injector.Backends(c).Plugin.Disable(c.Request.Context(), name); err != nil {
			c.Error(err)
			return
		}
	}

	if err := injector.Backends(c).Plugin.SetConfig(c.Request.Context(), name, req.Config); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
}
