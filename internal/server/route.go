package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	queryrouter "github.com/bluemir/wikinote/internal/query-router"
	"github.com/bluemir/wikinote/internal/server/middleware/auth"
	"github.com/bluemir/wikinote/internal/static"
)

var (
	authz = auth.Authz
)

func (server *Server) RegisterRoute(app gin.IRouter) {
	app.GET("/", server.redirectToFrontPage)

	{
		special := app.Group("/-")
		special.Group("/static", server.staticCache).StaticFS("/", static.Files.HTTPBox())

		special.GET("/auth/login", auth.Login)
		special.GET("/auth/profile", server.handler.Profile)

		// Register
		special.GET("/auth/register", server.static("/register.html"))
		special.POST("/auth/register", server.handler.Register)

		special.GET("/messages", server.handler.Messages)
		special.GET("/search", authz(Global, "search"), server.handler.Search)

		special.GET("/admin", authz(Global, "read"), server.handler.Admin)
	}
	{
		api := app.Group("/-/api")
		api.POST("/preview", server.handler.Preview) // render body
		api.GET("/me", server.handler.Me)
		api.GET("auth/can/:verb/*kind", server.handler.Can)
	}

	// plugins
	server.Backend.Plugin.RouteHook(app.Group("/~"))

	{
		// - GET            render file or render functional page
		//   - edit      : show editor
		//   - delete    : show delete check page
		//   - raw       : show raw text(not rendered)
		// - POST           create or update file with form submit
		// - PUT            create or update file with ajax
		// - DELETE         delete file

		pages := queryrouter.New()

		pages.GET("edit", authz(Page, "update"), server.handler.EditForm)
		pages.GET("raw", authz(Page, "read"), server.handler.Raw)
		pages.GET("delete", authz(Page, "delete"), server.static("/delete.html"))
		pages.GET("files", authz(Page, "update"), server.handler.Files)
		pages.GET("*", authz(Page, "read"), server.handler.View)
		pages.POST("*", authz(Page, "update"), server.handler.UpdateWithForm)
		pages.PUT("*", authz(Page, "update"), server.handler.Update)
		pages.DELETE("*", authz(Page, "delete"), server.handler.Delete)

		app.Use(rejectDotApp, pages.Handler)
	}
}
func (server *Server) redirectToFrontPage(c *gin.Context) {
	logrus.Debugf("redirect to front page: %s", server.frontPage)
	c.Redirect(http.StatusTemporaryRedirect, "/"+server.frontPage)
	c.Abort()
	return
}
func Page(c *gin.Context) (auth.Resource, error) {
	// get attributes
	return auth.KeyValues{
		"kind": "page",
		"path": c.Request.URL.Path,
	}, nil
}
func Global(c *gin.Context) (auth.Resource, error) {
	return auth.KeyValues{}, nil
}
func Admin(c *gin.Context) (auth.Resource, error) {
	return auth.KeyValues{
		"kind": "admin",
	}, nil
}
func rejectDotApp(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/.app") {
		c.Abort()
		return
	}
}
