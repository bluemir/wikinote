package server

import (
	"net/http"

	"github.com/bluemir/go-utils/auth"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkgs/backend"
	"github.com/bluemir/wikinote/pkgs/dist"
	"github.com/bluemir/wikinote/pkgs/renderer"
)

// gin context or session keys
const (
	BACKEND = "backend"
	SPECAIL = renderer.KEY_SPECIAL_ROUTE
	TOKEN   = renderer.KEY_TOKEN

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
		special.Use(BasicAuth)
		special.GET("/auth/login", HandleLogin)
		special.GET("/auth/logout", HandleLogout)

		special.POST("/api/preview", Action("edit"), HandlePreview) // render body
		special.GET("/search", Action("view", "search"), HandleSearch)

		special.GET("/user", Action("user"), HandleUserList)
		special.GET("/user/:id")
		special.PUT("/user")

		//special.GET("/api/users", Action("user"), HandleAPIUesrList)
		//special.GET("/api/users/:id", Action("user"), HandleAPIUesr)
		special.PUT("/api/users/:name/role", Action("user"), HandleAPIUserUpdateRole)
		// TODO make Action for API
	}

	app.Use(BasicAuth)
	app.NoRoute(func(c *gin.Context) {
		// GET            render file or render functional page
		// POST           create or update file with form submit
		// PUT            create or update file with ajax
		// DELETE         delete file

		switch c.Request.Method {
		case http.MethodGet:
			// check param
			// TODO auth
			if _, ok := c.GetQuery("edit"); ok {
				Do(c, HandleEditForm, "edit")
				return
			}
			if _, ok := c.GetQuery("attach"); ok {
				Do(c, HandleAttachForm, "attach")
				return
			}
			if _, ok := c.GetQuery("raw"); ok {
				Do(c, HandleRaw, "view")
				return
			}
			if _, ok := c.GetQuery("history"); ok {
				c.String(http.StatusNotImplemented, "text/plain", "not implemeted")
				return
			}
			if _, ok := c.GetQuery("backlinks"); ok {
				c.String(http.StatusNotImplemented, "text/plain", "not implemeted")
				return
			}
			if _, ok := c.GetQuery("move"); ok {
				c.String(http.StatusNotImplemented, "text/plain", "not implemeted")
				return
			}
			Do(c, HandleView, "view")

		case http.MethodPost:
			Do(c, HandleUpdateForm, "edit")
		case http.MethodPut:
			Do(c, HandleUpdate, "edit")
		case http.MethodDelete:
			//Do(c, HandleDelte, "edit")
			c.String(http.StatusNotImplemented, "text/plain", "not implemeted")
			return
		}
	})

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

func authCheck(c *gin.Context, actions ...string) int {
	// check user
	//      if not logined return 401
	//      if logined but not have permission return 403
	token, ok := c.Get(TOKEN)
	if ok {
		for _, action := range actions {
			if !Backend(c).Auth().Is(token).Allow(auth.Action(action)) {
				return http.StatusForbidden
			}
		}
	} else {
		for _, action := range actions {
			if !Backend(c).Auth().Is(backend.RoleGuest).Allow(auth.Action(action)) {
				logrus.Debugf("guest not allow %s", action)
				return http.StatusUnauthorized
			}
		}
	}
	return http.StatusOK
}

func Action(actions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch authCheck(c, actions...) {
		case http.StatusForbidden:
			// TODO make fordidden page
			c.HTML(http.StatusForbidden, "/errors/forbidden.html", renderer.Data{}.With(c))
			c.Abort()
			return
		case http.StatusUnauthorized:
			c.Header("WWW-Authenticate", AuthenicateString)
			c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
			c.Abort()
			return
		case http.StatusOK:
			return // pass
		}
	}
}
func Do(c *gin.Context, handler gin.HandlerFunc, actions ...string) {
	switch authCheck(c, actions...) {
	case http.StatusForbidden:
		// TODO make fordidden page
		c.HTML(http.StatusForbidden, "/errors/forbidden.html", renderer.Data{}.With(c))
		c.Abort()
		return
	case http.StatusUnauthorized:
		c.Header("WWW-Authenticate", AuthenicateString)
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	case http.StatusOK:
		handler(c)
		return // pass
	}
}
func Token(c *gin.Context) *auth.Token {
	token, ok := c.Get(TOKEN)
	if ok {
		return token.(*auth.Token)
	}
	return nil
}
