package server

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend"
	"github.com/bluemir/wikinote/backend/user"
	"github.com/bluemir/wikinote/server/renderer"
)

// gin context or session keys
const (
	BACKEND = "backend"
	ROLE    = "role"
	SPECAIL = "special-route"
	USER    = "user"
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
	store := sessions.NewCookieStore([]byte("__wikinote__"))
	app.Use(sessions.Sessions("session", store))

	app.Use(func(c *gin.Context) {
		c.Set(BACKEND, b)
	})

	app.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, b.Config().FrontPage)
	})
	special := app.Group("/!")
	{
		special.Use(func(c *gin.Context) {
			c.Set(SPECAIL, true)
		})
		special.StaticFS("/static/", rice.MustFindBox("../dist").HTTPBox())
		//auth
		special.GET("/auth/login", HandleLogin)
		special.GET("/auth/register", HandleRegisterForm)
		special.POST("/auth/register", HandleRegister)
		special.GET("/auth/logout", HandleLogout)

		special.Use(BasicAuth)
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
func User(c *gin.Context) *user.User {
	u, ok := c.Get(USER)
	if ok {
		return u.(*user.User)
	}
	return nil
}

func FlashMessage(c *gin.Context) renderer.MessageContext {
	return renderer.Of(c)
}

func authCheck(c *gin.Context, actions ...string) int {
	// check user
	//      if not logined return 401
	//      if logined but not have permission return 403
	u, ok := c.Get(USER)
	if ok {
		role := u.(*user.User).Role
		if Backend(c).Auth().IsAllow(role, actions...) {
			return http.StatusOK //pass
		} else {
			return http.StatusForbidden
		}
	} else {
		if Backend(c).Auth().IsAllow("guest", actions...) {
			return http.StatusOK //pass
		} else {
			return http.StatusUnauthorized
		}
	}
}

func Action(actions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		switch authCheck(c, actions...) {
		case http.StatusForbidden:
			// TODO make fordidden page
			c.HTML(http.StatusForbidden, "/errors/forbidden.html", renderer.Data{}.With(c))
			c.Abort()
		case http.StatusUnauthorized:
			c.Header("WWW-Authenticate", "Basic realm=\"Auth required!\"")
			c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
			c.Abort()
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
	case http.StatusUnauthorized:
		c.Header("WWW-Authenticate", "Basic realm=\"Auth required!\"")
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
	case http.StatusOK:
		handler(c)
		return // pass
	}
}
