package server

import (
	"net/http"

	"github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend"
	"github.com/bluemir/wikinote/server/renderer"
)

// gin context or session keys
const (
	BACKEND  = "backend"
	ROLE     = "role"
	USERNAME = "username"
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
	app.StaticFS("/!/static/", rice.MustFindBox("../dist").HTTPBox())
	//app.Group("/!/api")
	app.GET("/!/api/edit/*path")
	app.POST("/!/api/preview", Action("edit"), HandlePreview) // render body

	app.GET("/!/edit/*path", Action("edit"), HandleEditForm) // TODO must change base name because preview
	app.GET("/!/attach/*path", Action("edit"), HandleAttachForm)
	app.GET("/!/search", Action("view"), HandleSearch)
	app.GET("/!/history/*path")
	app.GET("/!/ws/*path")

	//auth
	app.GET("/!/auth/login", HandleLoginForm)
	app.POST("/!/auth/login", HandleLogin)
	app.GET("/!/auth/register", HandleRegisterForm)
	app.POST("/!/auth/register", HandleRegister)
	app.GET("/!/auth/logout", HandleLogout)

	app.GET("/!/user", Action("user"), HandleUserList)
	app.GET("/!/user/:id")
	app.PUT("/!/user")

	app.NoRoute(func(c *gin.Context) {
		// GET            render file
		// POST           create or update file with form submit
		// PUT            create or update file with ajax
		// DELETE         delete file

		switch c.Request.Method {
		case http.MethodGet:
			HandleView(c)
		case http.MethodPost:
			HandleUpdateForm(c)
		case http.MethodPut:
			HandleUpdate(c)
		case http.MethodDelete:
			HandleDelete(c)
		}
	})
	if len(conf.Domain) > 0 {
		logrus.Infof("Enable Auto TLS @ %s", conf.Domain)
		logrus.Infof("Run Server on %s", conf.Address)
		return autotls.Run(app, conf.Domain...)
	} else {
		logrus.Infof("Run Server on %s", conf.Address)
		return app.Run(conf.Address)
	}
}

func Backend(c *gin.Context) backend.Backend {
	return c.MustGet(BACKEND).(backend.Backend)
}
func Data(c *gin.Context) renderer.Data {
	return renderer.NewData(c)
}
func FlashMessage(c *gin.Context) renderer.MessageContext {
	return renderer.NewMessageContext(c)
}
func Action(actions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := Session(c).Role()
		for _, action := range actions {
			if Backend(c).Auth().IsAllow(role, action) {
				return //pass
			}
		}
		logrus.Debugf("403 role: %s, actions: %+v", role, actions)
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", Data(c))
		c.Abort()
	}
}
