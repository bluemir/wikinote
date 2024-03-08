package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/assets"
	queryrouter "github.com/bluemir/wikinote/internal/query-router"
	"github.com/bluemir/wikinote/internal/server/middleware/auth"
	"github.com/bluemir/wikinote/internal/server/middleware/auth/resource"
	"github.com/bluemir/wikinote/internal/server/middleware/auth/verb"
)

var (
	can = auth.Can
)

func (server *Server) route(app gin.IRouter, noRoute func(...gin.HandlerFunc)) {
	app.GET("/", server.redirectToFrontPage)

	{
		special := app.Group("/-")
		special.Group("/static", server.staticCache).StaticFS("/", http.FS(assets.Static))

		special.GET("/auth/login", server.handler.Login)
		special.GET("/auth/profile", server.handler.Profile)

		// Register
		special.GET("/auth/register", server.static("/register.html"))
		special.POST("/auth/register", server.handler.Register)

		special.GET("/messages", server.handler.Messages)
		special.GET("/search", can(verb.Search, resource.Global), server.handler.Search)

		special.GET("/admin", can(verb.Get, resource.AdminPage), server.handler.Admin)
	}
	{
		api := app.Group("/-/api", markAPI)
		api.POST("/preview", server.handler.Preview) // render body
		api.GET("/me", server.handler.Me)
		api.GET("auth/can/:verb/*kind", server.handler.Can)

		{
			v1 := api.Group("/v1")

			v1.GET("/me", server.handler.Me)
			v1.GET("/can/:verb/*kind", server.handler.Can)
		}

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

		pages.GET("edit", can(verb.Update, resource.Page), server.handler.EditForm)
		pages.GET("raw", can(verb.Get, resource.Page), server.handler.Raw)
		pages.GET("delete", can(verb.Delete, resource.Page), server.static("delete.html"))
		pages.GET("files", can(verb.Update, resource.Page), server.handler.Files)
		pages.GET("*", can(verb.Get, resource.Page), server.handler.View)
		pages.POST("*", can(verb.Update, resource.Page), server.handler.UpdateWithForm)
		pages.PUT("*", can(verb.Update, resource.Page), server.handler.Update)
		pages.DELETE("*", can(verb.Delete, resource.Page), server.handler.Delete)

		noRoute(rejectDotApp, markHTML, pages.Handler)
	}
}
func (server *Server) redirectToFrontPage(c *gin.Context) {
	logrus.Debugf("redirect to front page: %s", server.frontPage)
	c.Redirect(http.StatusTemporaryRedirect, "/"+server.frontPage)
	c.Abort()
}
func rejectDotApp(c *gin.Context) {
	if strings.HasPrefix(c.Request.URL.Path, "/.app") {
		c.Abort()
		return
	}
}
