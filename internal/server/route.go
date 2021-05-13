package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/query-router"
	"github.com/bluemir/wikinote/internal/static"
)

func (server *Server) RegisterRoute(app gin.IRouter) {
	redirectToFrontPage := func(c *gin.Context) {
		logrus.Debugf("redirect to front page: %s", server.Config.File.FrontPage)
		c.Redirect(http.StatusTemporaryRedirect, "/"+server.Config.File.FrontPage)
		c.Abort()
		return
	}
	app.GET("/", redirectToFrontPage)

	special := app.Group("/!")
	{
		special.Group("/static", staticCache).StaticFS("/", static.Files.HTTPBox())
		special.Group("/lib", staticCache).StaticFS("/", static.NodeModules.HTTPBox())

		// XXX for dev. must disable after dev
		// special.PUT("/api/users/:name/role", server.HandleAPIUserUpdateRole)
		special.GET("/login", server.Authn, redirectToFrontPage)

		// TODO user manager

		// Register
		special.GET("/auth/register", server.HandleRegisterForm)
		special.POST("/auth/register", server.HandleRegister)

		// auth
		special.Use(server.Authn)
		special.POST("/api/preview", server.HandlePreview) // render body
		special.GET("/search", server.Authz("search"), server.HandleSearch)

		// plugins
		server.Backend.Plugin.RouteHook(special.Group("/plugins"))
	}

	app.Use(server.Authn)
	// - GET            render file or render functional page
	//   - edit      : show editor
	//   - delete    : show delete check page
	//   - raw       : show raw text(not rendered)
	// - POST           create or update file with form submit
	// - PUT            create or update file with ajax
	// - DELETE         delete file

	queryRouter := queryrouter.New()
	queryRouter.Register(http.MethodGet, "edit", server.Authz("update"), server.HandleEditForm)
	queryRouter.Register(http.MethodGet, "raw", server.Authz("read"), server.HandleRaw)
	queryRouter.Register(http.MethodGet, "delete", server.Authz("delete"), server.HandleDeleteForm)
	queryRouter.Register(http.MethodGet, "attribute", server.Authz("read"), server.HandleAttributeGet)
	queryRouter.Register(http.MethodPut, "attribute", server.Authz("update"), server.HandleAttributeUpdate)
	queryRouter.Register(http.MethodGet, "*", server.Authz("read"), server.HandleView)
	queryRouter.Register(http.MethodPost, "*", server.Authz("update"), server.HandleUpdateWithForm)
	queryRouter.Register(http.MethodPut, "*", server.Authz("update"), server.HandleUpdate)
	queryRouter.Register(http.MethodDelete, "*", server.Authz("delete"), server.HandleDelete)

	app.Use(queryRouter.Handler)
}
