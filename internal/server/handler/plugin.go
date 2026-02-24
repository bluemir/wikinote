package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/server/injector"
)

func ListPlugins(c *gin.Context) error {
	plugins, err := injector.Backends(c).Plugin.List(c.Request.Context())
	if err != nil {
		return err
	}

	c.HTML(http.StatusOK, "admin/plugins/list.html", With(c, plugins))
	return nil
}
func GetPlugin(c *gin.Context) error {
	plugin, err := injector.Backends(c).Plugin.Get(c.Request.Context(), c.Param("name"))
	if err != nil {
		return err
	}

	c.HTML(http.StatusOK, "admin/plugins/config.html", With(c, plugin))
	return nil
}
func UpdatePlugin(c *gin.Context) error {
	name := c.Param("name")

	req := struct {
		IsEnabled bool   `form:"isEnabled"`
		Config    string `form:"config"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		return err
	}

	logrus.Tracef("name: %s, request: %+v", name, req)

	if err := injector.Backends(c).Plugin.SetConfig(c.Request.Context(), name, req.Config); err != nil {
		return err
	}

	if req.IsEnabled {
		if err := injector.Backends(c).Plugin.Enable(c.Request.Context(), name); err != nil {
			return err
		}
	} else {
		if err := injector.Backends(c).Plugin.Disable(c.Request.Context(), name); err != nil {
			return err
		}
	}

	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
	return nil
}
