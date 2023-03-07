package server

import (
	"net/http"
	"os"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/acme/autocert"

	"github.com/bluemir/wikinote/internal/backend"
	"github.com/bluemir/wikinote/internal/server/handler"
	authMiddleware "github.com/bluemir/wikinote/internal/server/middleware/auth"
)

type Config struct {
	Bind         string
	FrontPage    string
	EnableHttps  bool
	HttpsDomain  string
	AutoTLSCache string
}

func NewConfig() *Config {
	return &Config{}
}

type Server struct {
	*backend.Backend
	handler   *handler.Handler
	frontPage string
	etag      string
}

func Run(b *backend.Backend, conf *Config) error {
	h, err := handler.New(b)
	if err != nil {
		return err
	}

	server := &Server{
		Backend:   b,
		handler:   h,
		frontPage: conf.FrontPage,
	}

	app := gin.New()

	// Log
	writer := logrus.WithFields(logrus.Fields{}).WriterLevel(logrus.DebugLevel)
	defer writer.Close()
	app.Use(gin.LoggerWithWriter(writer))

	// Recovery
	//app.Use(gin.Recovery())
	app.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, err)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	app.Use(location.Default(), fixURL)

	// Renderer
	if r, err := NewRender(); err != nil {
		return errors.WithStack(err)
	} else {
		app.HTMLRender = r
	}

	app.Use(authMiddleware.Middleware(server.Backend.Auth))

	app.GET("/favicon.ico", NotFound)
	// Register Routing
	server.RegisterRoute(app)

	cacheDir := conf.AutoTLSCache
	if cacheDir == "" {
		cacheDir = b.ConfigPath("cert-cache")
	}
	os.MkdirAll(cacheDir, 0700)

	if conf.EnableHttps {
		logrus.Infof("Run Server with AutoTLS")
		return autotls.RunWithManager(app, &autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(conf.HttpsDomain),
			Cache:      autocert.DirCache(cacheDir),
		})
	} else {
		logrus.Infof("Run Server on %s", conf.Bind)
		return app.Run(conf.Bind)
	}
}
func NotFound(c *gin.Context) {
	c.String(http.StatusNotFound, "Not Found")
}
