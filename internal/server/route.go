package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/assets"
	"github.com/bluemir/wikinote/internal/plugins"
	queryrouter "github.com/bluemir/wikinote/internal/query-router"
	"github.com/bluemir/wikinote/internal/server/handler"
	"github.com/bluemir/wikinote/internal/server/handler/auth/resource"
	"github.com/bluemir/wikinote/internal/server/handler/auth/verb"
	"github.com/bluemir/wikinote/internal/server/middleware/cache"
)

var (
	can = handler.Can
)

func w(fn func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fn(c); err != nil {
			c.Error(err)
			c.Abort()
		}
	}
}

func (server *Server) route(app gin.IRouter, noRoute func(...gin.HandlerFunc), plugins *plugins.Manager) {
	// favicon
	app.GET("/favicon.ico", w(handler.NotFound))

	app.GET("/", server.redirectToFrontPage)

	// static
	app.Group("/-/static/", staticCache).Group(cache.Rev()).StaticFS("/", http.FS(assets.Static))
	{
		// APIs
		{
			v1 := app.Group("/-/api/v1", markAPI)

			//v1.GET("/can/:verb/*kind", handler.CanAPI)
			//v1.GET("/me", api(handler.Me))

			v1.GET("/preview", w(handler.Preview)) // render body
		}
	}
	{
		// system pages
		system := app.Group("/-", markHTML)

		system.GET("/welcome", html("system/welcome.html"))

		system.GET("/auth/login", html("system/user/login.html"))
		system.POST("/auth/login", w(handler.Login))
		system.GET("/auth/logout", w(handler.Logout))
		system.GET("/auth/profile", w(handler.Profile))

		system.GET("/auth/register", html("system/user/register.html"))
		system.POST("/auth/register", w(handler.Register))

		system.GET("/search", can(verb.Search, resource.Page), w(handler.Search))

		system.GET("/admin", can(verb.Get, resource.AdminPage), html("admin/index.html"))

		// common
		system.GET("/admin/common ", can(verb.Get, resource.AdminPage), html("admin/common.html"))
		system.POST("/admin", can(verb.Update, resource.AdminPage), w(handler.UpdateCommonSetting))

		// iam > users
		system.GET("/admin/iam/users", can(verb.List, resource.Users), w(handler.ListUsers))
		system.GET("/admin/iam/users/:username", can(verb.Get, resource.Users), w(handler.GetUser))
		system.POST("/admin/iam/users/:username", can(verb.Get, resource.Users), w(handler.UpdateUser))

		// iam > groups
		system.POST("/admin/iam/groups", can(verb.Create, resource.Users), w(handler.CreateGroup))
		system.GET("/admin/iam/groups", can(verb.List, resource.Users), w(handler.ListGroups))
		system.GET("/admin/iam/groups/:groupName", can(verb.Get, resource.Groups), w(handler.GetGroup))
		system.DELETE("/admin/iam/groups/:groupName", can(verb.Delete, resource.Groups), w(handler.DeleteGroup))

		// iam > role
		system.POST("/admin/iam/roles", can(verb.Create, resource.Roles), w(handler.CreateRole))
		system.GET("/admin/iam/roles", can(verb.List, resource.Roles), w(handler.ListRoles))
		system.GET("/admin/iam/roles/:roleName", can(verb.Get, resource.Roles), w(handler.GetRole))
		system.POST("/admin/iam/roles/:roleName", can(verb.Update, resource.Roles), w(handler.UpdateRole))
		system.DELETE("/admin/iam/roles/:roleName", can(verb.Delete, resource.Roles), w(handler.DeleteRole))

		// iam > assigns
		system.POST("/admin/iam/assigns", can(verb.Create, resource.Assigns), w(handler.CreateAssign))
		system.GET("/admin/iam/assigns", can(verb.List, resource.Assigns), w(handler.ListAssigns))
		system.GET("/admin/iam/assigns/:subjectKind", can(verb.List, resource.Assigns), w(handler.GetAssign))
		system.PUT("/admin/iam/assigns/:subjectKind", can(verb.Update, resource.Assigns), w(handler.UpdateAssign))
		system.DELETE("/admin/iam/assigns/:subjectKind", can(verb.Delete, resource.Assigns), w(handler.DeleteAssign))
		system.GET("/admin/iam/assigns/:subjectKind/:subjectName", can(verb.List, resource.Assigns), w(handler.GetAssign))
		system.PUT("/admin/iam/assigns/:subjectKind/:subjectName", can(verb.Update, resource.Assigns), w(handler.UpdateAssign))
		system.DELETE("/admin/iam/assigns/:subjectKind/:subjectName", can(verb.Delete, resource.Assigns), w(handler.DeleteAssign))

		// plugins
		system.GET("/admin/plugins", can(verb.List, resource.Plugins), w(handler.ListPlugins))
		system.GET("/admin/plugins/:name", can(verb.Get, resource.Plugins), w(handler.GetPlugin))
		system.PUT("/admin/plugins/:name", can(verb.Update, resource.Plugins), w(handler.UpdatePlugin))

		system.GET("/admin/messages", can(verb.List, resource.Messages), w(handler.ListAllMessages))

		system.GET("/admin/metadata", can(verb.Get, resource.Metadata), w(handler.ListMetadata))

		system.GET("/initialize", w(handler.RequestInitialize))
		system.GET("/initialize/:code", w(handler.Initialze))
		system.POST("/initialize/:code", w(handler.InitialzeAccept))
	}

	// plugins
	//plugins.RouteHook(app.Group("/~"))
	app.Any("/~/:name/*path", plugins.HandleHTTPRequest)

	// reject url
	app.Any("/.app/*path", w(handler.NotFound))
	app.Any("/.git/*path", w(handler.NotFound))
	app.Use(IfHasPrefix("/-/", w(handler.NotFound)))
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

		pages.GET("edit", can(verb.Update, resource.Page), w(handler.EditForm))
		pages.GET("move", can(verb.Update, resource.Page), html("notes/move.html"))
		pages.POST("move", can(verb.Update, resource.Page), w(handler.MoveNote))
		pages.GET("raw", can(verb.Get, resource.Page), w(handler.Raw))
		pages.GET("delete", can(verb.Delete, resource.Page), html("notes/delete.html"))
		pages.POST("file", can(verb.Update, resource.Page), w(handler.UploadFiles))        // upload file, with multipart encoding, create or overwrite file
		pages.PUT("file", can(verb.Update, resource.Page), w(handler.UploadFileToReplace)) // upload file. with multipart encoding. replace file
		pages.GET("*", can(verb.Get, resource.Page), w(handler.View))
		pages.PUT("*", can(verb.Update, resource.Page), w(handler.UpdateWithForm))
		pages.DELETE("*", can(verb.Delete, resource.Page), w(handler.DeleteNote))

		//app.Any("/*path", pages.Handler)
		noRoute(pages.Handler)
	}

}
func IfHasPrefix(prefix string, fns ...gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, prefix) {
			for _, fn := range fns {
				if c.IsAborted() {
					return
				}
				fn(c)
			}
		}
	}
}
func (server *Server) redirectToFrontPage(c *gin.Context) {
	logrus.Debugf("redirect to front page: %s", server.frontPage)
	c.Redirect(http.StatusTemporaryRedirect, "/"+server.frontPage)
	c.Abort()
}
