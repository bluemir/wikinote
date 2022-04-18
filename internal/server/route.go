package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/auth"
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
		special.GET("/search", server.Authz(Global, "search"), server.handler.Search)

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
		pages.GET("edit", server.Authz(Page, "update"), server.handler.EditForm)
		pages.GET("raw", server.Authz(Page, "read"), server.handler.Raw)
		pages.GET("delete", server.Authz(Page, "delete"), server.handler.DeleteForm)
		pages.GET("attribute", server.Authz(PageAttr, "read"), server.handler.AttributeGet)
		pages.PUT("attribute", server.Authz(PageAttr, "update"), server.handler.AttributeUpdate)
		pages.GET("*", server.Authz(Page, "read"), server.handler.View)
		pages.POST("*", server.Authz(Page, "update"), server.handler.UpdateWithForm)
		pages.PUT("*", server.Authz(Page, "update"), server.handler.Update)
		pages.DELETE("*", server.Authz(Page, "delete"), server.handler.Delete)
	}
	app.Use(pages.Handler)
}
func (server *Server) redirectToFrontPage(c *gin.Context) {
	logrus.Debugf("redirect to front page: %s", server.Config.File.FrontPage)
	c.Redirect(http.StatusTemporaryRedirect, "/"+server.Config.File.FrontPage)
	c.Abort()
	return
}
func Page(c *gin.Context) (auth.Resource, error) {
	return auth.KeyValues{}, nil
}
func PageAttr(c *gin.Context) (auth.Resource, error) {
	return auth.KeyValues{}, nil
}
func Global(c *gin.Context) (auth.Resource, error) {
	return auth.KeyValues{}, nil
}
