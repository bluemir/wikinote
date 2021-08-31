package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/query-router"
	"github.com/bluemir/wikinote/internal/static"
)

func (server *Server) RegisterRoute(app gin.IRouter) {
	app.GET("/", server.redirectToFrontPage)

	special := app.Group("/!")
	{
		special.Group("/static", server.staticCache).StaticFS("/", static.Files.HTTPBox())

		// XXX for dev. must disable after dev
		// special.PUT("/api/users/:name/role", server.HandleAPIUserUpdateRole)
		special.GET("/login", server.Authn, server.redirectToFrontPage)

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

	pages := queryrouter.New()
	{
		pages.GET("edit", server.Authz("update"), server.HandleEditForm)
		pages.GET("raw", server.Authz("read"), server.HandleRaw)
		pages.GET("delete", server.Authz("delete"), server.HandleDeleteForm)
		pages.GET("attribute", server.Authz("read"), server.HandleAttributeGet)
		pages.PUT("attribute", server.Authz("update"), server.HandleAttributeUpdate)
		pages.GET("*", server.Authz("read"), server.HandleView)
		pages.POST("*", server.Authz("update"), server.HandleUpdateWithForm)
		pages.PUT("*", server.Authz("update"), server.HandleUpdate)
		pages.DELETE("*", server.Authz("delete"), server.HandleDelete)
	}
	app.Use(pages.Handler)
}
func (server *Server) redirectToFrontPage(c *gin.Context) {
	logrus.Debugf("redirect to front page: %s", server.Config.File.FrontPage)
	c.Redirect(http.StatusTemporaryRedirect, "/"+server.Config.File.FrontPage)
	c.Abort()
	return
}
