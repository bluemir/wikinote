package server

import (
	"net/http"
	"strings"

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
	// favicon
	app.GET("/favicon.ico", handler.NotFound)

	app.GET("/", server.redirectToFrontPage)

	{
		// APIs
		api := app.Group("/-/api", markAPI)
		{
			v1 := api.Group("/v1")

			//v1.GET("/can/:verb/*kind", handler.CanAPI)
			//v1.GET("/me", handler.Me)

			v1.GET("/preview", handler.Preview) // render body
		}
	}
	{
		// system pages
		system := app.Group("/-", markHTML)
		system.Group("/static", staticCache).StaticFS("/", http.FS(assets.Static))

		system.GET("/welcome", html("system/welcome.html"))

		system.GET("/auth/login", html("system/user/login.html"))
		system.POST("/auth/login", handler.Login)
		system.GET("/auth/logout", handler.Logout)
		system.GET("/auth/profile", handler.Profile)

		system.GET("/auth/register", html("system/user/register.html"))
		system.POST("/auth/register", handler.Register)

		system.GET("/search", can(verb.Search, resource.Page), handler.Search)

		system.GET("/admin", can(verb.Get, resource.AdminPage), html("admin/index.html"))
		system.GET("/admin/iam/users", can(verb.List, resource.Users), handler.ListUsers)
		system.GET("/admin/iam/users/:username", can(verb.Get, resource.Users), handler.GetUser)
		system.POST("/admin/iam/users/:username", can(verb.Get, resource.Users), handler.UpdateUser)

		// iam > groups
		system.POST("/admin/iam/groups", can(verb.Create, resource.Users), handler.CreateGroup)
		system.GET("/admin/iam/groups", can(verb.List, resource.Users), handler.ListGroups)
		system.GET("/admin/iam/groups/:groupName", can(verb.Get, resource.Groups), handler.GetGroup)
		system.DELETE("/admin/iam/groups/:groupName", can(verb.Delete, resource.Groups), handler.DeleteGroup)

		// iam > role
		system.POST("/admin/iam/roles", can(verb.Create, resource.Roles), handler.CreateRole)
		system.GET("/admin/iam/roles", can(verb.List, resource.Roles), handler.ListRoles)
		system.GET("/admin/iam/roles/:roleName", can(verb.Get, resource.Roles), handler.GetRole)
		system.POST("/admin/iam/roles/:roleName", can(verb.Update, resource.Roles), handler.UpdateRole)
		system.DELETE("/admin/iam/roles/:roleName", can(verb.Delete, resource.Roles), handler.DeleteRole)

		// iam > assigns
		system.POST("/admin/iam/assigns", can(verb.Create, resource.Assigns), handler.CreateAssign)
		system.GET("/admin/iam/assigns", can(verb.List, resource.Assigns), handler.ListAssigns)
		system.GET("/admin/iam/assigns/:subjectKind", can(verb.List, resource.Assigns), handler.GetAssign)
		system.PUT("/admin/iam/assigns/:subjectKind", can(verb.Update, resource.Assigns), handler.UpdateAssign)
		system.DELETE("/admin/iam/assigns/:subjectKind", can(verb.Delete, resource.Assigns), handler.DeleteAssign)
		system.GET("/admin/iam/assigns/:subjectKind/:subjectName", can(verb.List, resource.Assigns), handler.GetAssign)
		system.PUT("/admin/iam/assigns/:subjectKind/:subjectName", can(verb.Update, resource.Assigns), handler.UpdateAssign)
		system.DELETE("/admin/iam/assigns/:subjectKind/:subjectName", can(verb.Delete, resource.Assigns), handler.DeleteAssign)

		system.GET("/admin/messages", can(verb.List, resource.Messages), handler.ListAllMessages)

		system.GET("/initialize", handler.RequestInitialize)
		system.GET("/initialize/:code", handler.Initialze)
		system.POST("/initialize/:code", handler.InitialzeAccept)

	}

	// plugins
	plugins.RouteHook(app.Group("/~"))

	// reject url
	app.Any("/.app/*path", handler.NotFound)
	app.Any("/.git/*path", handler.NotFound)
	app.Use(IfHasPrefix("/-/", handler.NotFound))
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
		pages.GET("move", can(verb.Update, resource.Page), html("notes/move.html"))
		pages.POST("move", can(verb.Update, resource.Page), handler.MoveNote)
		pages.GET("raw", can(verb.Get, resource.Page), handler.Raw)
		pages.GET("delete", can(verb.Delete, resource.Page), html("notes/delete.html"))
		pages.POST("file", can(verb.Update, resource.Page), handler.UploadFiles)        // upload file, with multipart encoding, create or overwrite file
		pages.PUT("file", can(verb.Update, resource.Page), handler.UploadFileToReplace) // upload file. with multipart encoding. replace file
		pages.GET("*", can(verb.Get, resource.Page), handler.View)
		pages.PUT("*", can(verb.Update, resource.Page), handler.UpdateWithForm)
		pages.DELETE("*", can(verb.Delete, resource.Page), handler.DeleteNote)

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
