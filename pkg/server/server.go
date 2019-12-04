package server

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkg/backend"
	"github.com/bluemir/wikinote/pkg/dist"
	"github.com/bluemir/wikinote/pkg/query-router"
)

type Config struct {
	Bind       string
	TLSDomains []string
}

func Run(b *backend.Backend, conf *Config) error {
	server := &Server{b}

	app := gin.New()

	// Log
	writer := logrus.New().Writer()
	defer writer.Close()
	app.Use(gin.LoggerWithWriter(writer))

	// Recovery
	app.Use(gin.Recovery())

	// Session
	store := cookie.NewStore([]byte("__wikinote__"))
	app.Use(sessions.Sessions("session", store))

	// Renderer
	app.HTMLRender = loadTemplates()

	// Root
	indexRouter := queryrouter.New()
	redirectToFrontPage := func(c *gin.Context) {
		logrus.Debugf("redirect to front page: %s", b.Config.File.FrontPage)
		c.Redirect(http.StatusTemporaryRedirect, b.Config.File.FrontPage)
		c.Abort()
		return
	}
	indexRouter.Register(http.MethodGet, "login", server.Authn, redirectToFrontPage)
	indexRouter.Register(http.MethodGet, "*", redirectToFrontPage)
	app.GET("/", indexRouter.Handler)

	special := app.Group("/!")
	{
		special.StaticFS("/static/", dist.Files.HTTPBox())

		// XXX for dev. must disable after dev
		//special.PUT("/api/users/:name/role", server.HandleAPIUserUpdateRole)

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
	queryRouter.Register(http.MethodGet, "*", server.Authz("read"), server.HandleView)
	queryRouter.Register(http.MethodPost, "*", server.Authz("update"), server.HandleUpdateWithForm)
	queryRouter.Register(http.MethodPut, "*", server.Authz("update"), server.HandleUpdate)
	queryRouter.Register(http.MethodDelete, "*", server.Authz("delete"), server.HandleDelete)

	app.NoRoute(queryRouter.Handler)

	if len(conf.TLSDomains) > 0 {
		logrus.Warn("ignore bind or port option")
		logrus.Info("Run http redirect server")
		go http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logrus.Infof("hit the http request: %s", r.RequestURI)
			http.Redirect(w, r, "https://"+conf.TLSDomains[0]+r.RequestURI, http.StatusPermanentRedirect)
		}))

		logrus.Infof("Enable Auto TLS @ %s", conf.TLSDomains)
		logrus.Infof("Run Server")

		return autotls.Run(app, conf.TLSDomains...)
	} else {
		logrus.Infof("Run Server on %s", conf.Bind)
		return app.Run(conf.Bind)
	}
}

type Server struct {
	*backend.Backend
}

func (server *Server) HandleNotImplemented(c *gin.Context) {
	c.String(http.StatusNotImplemented, "text/plain", "not implemeted")
	c.Abort()
	return
}
