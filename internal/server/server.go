package server

import (
	"net/http"

	"github.com/gin-contrib/location"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/backend"
	"github.com/bluemir/wikinote/internal/server/handler"
	authMiddleware "github.com/bluemir/wikinote/internal/server/middleware/auth"
)

type Config struct {
	Bind string
}

func NewConfig() *Config {
	return &Config{}
}

type Server struct {
	*backend.Backend
	handler *handler.Handler
	etag    string
}

func Run(b *backend.Backend, conf *Config) error {
	h, err := handler.New(b)
	if err != nil {
		return err
	}

	server := &Server{
		Backend: b,
		handler: h,
	}

	app := gin.New()

	// Log
	writer := logrus.WithFields(logrus.Fields{}).WriterLevel(logrus.DebugLevel)
	defer writer.Close()
	app.Use(gin.LoggerWithWriter(writer))

	// Recovery
	app.Use(gin.Recovery())

	app.Use(location.Default(), fixURL)

	// Session
	store := cookie.NewStore([]byte("__wikinote__"))
	app.Use(sessions.Sessions("session", store))

	// Renderer
	if html, err := NewRenderer(); err != nil {
		return errors.WithStack(err)
	} else {
		app.SetHTMLTemplate(html)
	}

	app.Use(authMiddleware.Middleware(server.Backend.Auth))

	app.GET("/favicon.ico", NotFound)
	// Register Routing
	server.RegisterRoute(app)

	logrus.Infof("Run Server on %s", conf.Bind)
	return app.Run(conf.Bind)
}
func NotFound(c *gin.Context) {
	c.String(http.StatusNotFound, "Not Found")
}
