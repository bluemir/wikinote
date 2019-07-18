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

func Run(b backend.Backend, conf *Config) error {
	app := gin.New()
	writer := logrus.New().Writer()
	defer writer.Close()
	app.Use(gin.LoggerWithWriter(writer))
	app.Use(gin.Recovery())

	app.HTMLRender = renderer.NewRenderer()

	// TODO rendom string or config
	store := cookie.NewStore([]byte("__wikinote__"))
	app.Use(sessions.Sessions("session", store))

	app.Use(func(c *gin.Context) {
		c.Set(BACKEND, b)
	})

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
		special.GET("/auth/register", HandleRegisterForm)
		special.POST("/auth/register", HandleRegister)

		//auth
		special.Use(BasicAuthn)
		special.GET("/auth/login", HandleLogin)
		special.GET("/auth/logout", HandleLogout)

		special.POST("/api/preview", Authz("edit"), HandlePreview) // render body
		special.GET("/search", Authz("search"), HandleSearch)

		special.GET("/user", Authz("user"), HandleUserList)
		special.GET("/user/:id")
		special.PUT("/user")

		//special.GET("/api/users", Action("user"), HandleAPIUesrList)
		//special.GET("/api/users/:id", Action("user"), HandleAPIUesr)
		special.PUT("/api/users/:name/role", Authz("user"), HandleAPIUserUpdateRole)
		// TODO make Action for API
	}

	app.Use(BasicAuthn)

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
	queryRouter.Register(http.MethodGet, "edit", Authz("edit"), HandleEditForm)
	queryRouter.Register(http.MethodGet, "attach", Authz("attach"), HandleAttachForm)
	queryRouter.Register(http.MethodGet, "raw", Authz("raw"), HandleRaw)
	queryRouter.Register(http.MethodGet, "history", HandleNotImplemented)
	queryRouter.Register(http.MethodGet, "backlinks", HandleNotImplemented)
	queryRouter.Register(http.MethodGet, "move", HandleNotImplemented)
	queryRouter.Register(http.MethodGet, "*", HandleView)
	queryRouter.Register(http.MethodPost, "*", Authz("edit"), HandleUpdateForm)
	queryRouter.Register(http.MethodPut, "*", Authz("edit"), HandleUpdate)
	queryRouter.Register(http.MethodDelete, "*", HandleNotImplemented)

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

func Backend(c *gin.Context) backend.Backend {
	return c.MustGet(BACKEND).(backend.Backend)
}

func FlashMessage(c *gin.Context) renderer.MessageContext {
	return renderer.Of(c)
}

func HandleNotImplemented(c *gin.Context) {
	c.String(http.StatusNotImplemented, "text/plain", "not implemeted")
	c.Abort()
	return
}
