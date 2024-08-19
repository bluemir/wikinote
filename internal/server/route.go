package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/assets"
	"github.com/bluemir/wikinote/internal/plugins"
	queryrouter "github.com/bluemir/wikinote/internal/query-router"
	"github.com/bluemir/wikinote/internal/server/handler"
	"github.com/bluemir/wikinote/internal/server/handler/auth/resource"
	"github.com/bluemir/wikinote/internal/server/handler/auth/verb"
)

var (
	can = handler.Can
)

func (server *Server) route(app gin.IRouter, noRoute func(...gin.HandlerFunc), plugins *plugins.Manager) {
	app.GET("/", server.redirectToFrontPage)

	{
		// APIs
		api := app.Group("/-/api", markAPI)

		{
			v1 := api.Group("/v1")

			v1.GET("/can/:verb/*kind", handler.CanAPI)
			v1.GET("/me", handler.Me)

			v1.GET("/preview", handler.Preview) // render body

			//v1.POST("/users", handler.Register)

			//v1.GET("/iam/users", handler.ListUsers)
			//v1.GET("/iam/groups", handler.ListGroups)
			//v1.GET("/iam/roles", handler.ListRoles)

			//v1.GET("/config")

			//v1.GET("/events", handler.StreamEvents)
		}
	}
	{
		// system pages
		system := app.Group("/-", markHTML)
		system.Group("/static", staticCache).StaticFS("/", http.FS(assets.Static))

		system.GET("/welcome", html("welcome.html"))

		system.GET("/auth/login", html("login.html"))
		system.POST("/auth/login", handler.Login)
		system.GET("/auth/profile", handler.Profile)

		system.GET("/auth/register", html("register.html"))
		system.POST("/auth/register", handler.Register)

		system.GET("/messages", handler.Messages)
		system.GET("/search", can(verb.Search, resource.Global), handler.Search)

		system.GET("/admin", can(verb.Get, resource.AdminPage), html("admin/index.html"))
		system.GET("/admin/iam/users", can(verb.List, resource.Users), handler.ListUsers)
		system.GET("/admin/iam/groups", can(verb.List, resource.Users), handler.ListGroups)
		system.GET("/admin/iam/roles", can(verb.List, resource.Roles), handler.ListRoles)
	}

	// plugins
	plugins.RouteHook(app.Group("/~"))

	app.GET("/.app/*path", notFound)
	{
		// normal pages
		// - GET            render file or render functional page
		//   - edit      : show editor
		//   - delete    : show delete check page
		//   - raw       : show raw text(not rendered)
		// - POST           create or update file with form submit
		// - PUT            create or update file with ajax
		// - DELETE         delete file

		pages := queryrouter.New()

		pages.GET("edit", can(verb.Update, resource.Page), handler.EditForm)
		pages.GET("raw", can(verb.Get, resource.Page), handler.Raw)
		pages.GET("delete", can(verb.Delete, resource.Page), html("delete.html"))
		pages.GET("files", can(verb.Update, resource.Page), handler.Files)
		pages.GET("*", can(verb.Get, resource.Page), handler.View)
		pages.POST("*", can(verb.Update, resource.Page), handler.UpdateWithForm)
		pages.PUT("*", can(verb.Update, resource.Page), handler.Update)
		pages.DELETE("*", can(verb.Delete, resource.Page), handler.Delete)

		//app.Any("/*path", pages.Handler)
		noRoute(pages.Handler)
	}

}
func (server *Server) redirectToFrontPage(c *gin.Context) {
	logrus.Debugf("redirect to front page: %s", server.frontPage)
	c.Redirect(http.StatusTemporaryRedirect, "/"+server.frontPage)
	c.Abort()
}

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}
func (e NotFoundError) Code() int {
	return http.StatusNotFound
}
func notFound(c *gin.Context) {
	// TODO return not found error
	c.Error(NotFoundError(c.FullPath() + " not found"))
	c.Abort()
}
