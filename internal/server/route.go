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
		special.POST("/api/preview", server.handler.Preview) // render body
		special.GET("/search", server.Authz("search"), server.handler.Search)

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
		pages.GET("edit", server.Authz("update"), server.handler.EditForm)
		pages.GET("raw", server.Authz("read"), server.handler.Raw)
		pages.GET("delete", server.Authz("delete"), server.handler.DeleteForm)
		pages.GET("attribute", server.Authz("read"), server.handler.AttributeGet)
		pages.PUT("attribute", server.Authz("update"), server.handler.AttributeUpdate)
		pages.GET("*", server.Authz("read"), server.handler.View)
		pages.POST("*", server.Authz("update"), server.handler.UpdateWithForm)
		pages.PUT("*", server.Authz("update"), server.handler.Update)
		pages.DELETE("*", server.Authz("delete"), server.handler.Delete)
	}
	app.Use(pages.Handler)
}
func (server *Server) redirectToFrontPage(c *gin.Context) {
	logrus.Debugf("redirect to front page: %s", server.Config.File.FrontPage)
	c.Redirect(http.StatusTemporaryRedirect, "/"+server.Config.File.FrontPage)
	c.Abort()
	return
}
