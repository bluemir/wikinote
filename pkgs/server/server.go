package server

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkgs/backend"
	"github.com/bluemir/wikinote/pkgs/dist"
	"github.com/bluemir/wikinote/pkgs/query-router"
	"github.com/bluemir/wikinote/pkgs/renderer"
)

// gin context or session keys
const (
	BACKEND      = "backend"
	SPECAIL      = renderer.KEY_SPECIAL_ROUTE
	TOKEN        = renderer.KEY_TOKEN
	AUTH_CONTEXT = "auth-context"

	Relam             = "Wikinote"
	AuthenicateString = `Basic realm="` + Relam + `"`
)

type Config struct {
	Address string
	Domain  []string
}
type Server struct {
	backend.Backend
}

func Run(b backend.Backend, conf *Config) error {
	server := &Server{b}

	app := gin.New()
	writer := logrus.New().Writer()
	defer writer.Close()

	app.Use(gin.LoggerWithWriter(writer))
	app.Use(gin.Recovery())

	app.HTMLRender = renderer.NewRenderer()

	// TODO rendom string or config
	store := cookie.NewStore([]byte("__wikinote__"))
	app.Use(sessions.Sessions("session", store))

	b.Plugin().RegisterRouter(app.Group("/!/plugins"))

	app.GET("/", func(c *gin.Context) {
		logrus.Debugf("redirect to front page: %s", b.Config().FrontPage)
		c.Redirect(http.StatusTemporaryRedirect, b.Config().FrontPage)
	})
	special := app.Group("/!")
	{
		special.Use(func(c *gin.Context) {
			c.Set(SPECAIL, true)
		})
		special.StaticFS("/static/", dist.Apps.HTTPBox())

		// register
		special.GET("/auth/register", server.HandleRegisterForm)
		special.POST("/auth/register", server.HandleRegister)

		//auth
		special.Use(server.BasicAuthn)
		special.GET("/auth/login", server.HandleLogin)
		special.GET("/auth/logout", server.HandleLogout)

		special.POST("/api/preview", server.Authz("preview"), server.HandlePreview) // render body
		special.GET("/search", server.Authz("search"), server.HandleSearch)

		special.GET("/user", server.Authz("user"), server.HandleUserList)
		special.GET("/user/:id")
		special.PUT("/user")

		//special.GET("/api/users", Action("user"), HandleAPIUesrList)
		//special.GET("/api/users/:id", Action("user"), HandleAPIUesr)
		special.PUT("/api/users/:name/role", server.Authz("user"), server.HandleAPIUserUpdateRole)
		special.PUT("/api/users/:name/attr/*key", server.Authz("user"), server.HandleAPIPutUserAttr)
		// curl -XPUT 'http://localhost:4000/!/api/users/{{user}}/role' -u "root:{{rootKey}}" --data "{{role}}"
		// TODO make Action for API
	}

	app.Use(server.BasicAuthn)

	// - GET            render file or render functional page
	//   - edit      : show editor
	//   - attach    : show attachment and upload form
	//   - raw       : show raw text(not rendered)
	//   - history   : show histroy
	//   - backlinks : show backlinks
	//   - move      : show move confirm window
	// - POST           create or update file with form submit
	// - PUT            create or update file with ajax
	// - DELETE         delete file

	queryRouter := queryrouter.New()
	queryRouter.Register(http.MethodGet, "edit", server.Authz("edit"), server.HandleEditForm)
	queryRouter.Register(http.MethodGet, "attach", server.Authz("attach"), server.HandleAttachForm)
	queryRouter.Register(http.MethodGet, "raw", server.Authz("raw"), server.HandleRaw)
	queryRouter.Register(http.MethodGet, "history", server.HandleNotImplemented)
	queryRouter.Register(http.MethodGet, "backlinks", server.HandleNotImplemented)
	queryRouter.Register(http.MethodGet, "move", server.HandleNotImplemented)
	queryRouter.Register(http.MethodGet, "*", server.Authz("view"), server.HandleView)
	queryRouter.Register(http.MethodPost, "*", server.Authz("update"), server.HandleUpdateForm)
	queryRouter.Register(http.MethodPut, "*", server.Authz("update"), server.HandleUpdate)
	queryRouter.Register(http.MethodDelete, "*", server.HandleNotImplemented)

	server.Plugin().RegisterAction(queryRouter, server.Authz)

	app.NoRoute(queryRouter.Handler)

	if len(conf.Domain) > 0 {
		logrus.Warn("ignore bind or port option")
		logrus.Info("Run http redirect server")
		go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logrus.Infof("hit the http request: %s", r.RequestURI)
			http.Redirect(w, r, "https://"+conf.Domain[0]+r.RequestURI, http.StatusPermanentRedirect)
		}))

		logrus.Infof("Enable Auto TLS @ %s", conf.Domain)
		logrus.Infof("Run Server")
		return autotls.Run(app, conf.Domain...)
	} else {
		logrus.Infof("Run Server on %s", conf.Address)
		return app.Run(conf.Address)
	}
}

func FlashMessage(c *gin.Context) renderer.MessageContext {
	return renderer.Of(c)
}

func (server *Server) HandleNotImplemented(c *gin.Context) {
	c.String(http.StatusNotImplemented, "text/plain", "not implemeted")
	c.Abort()
	return
}
