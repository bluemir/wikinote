package server

import (
	"context"
	"net/http"

	"github.com/gin-contrib/location"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/backend"
	"github.com/bluemir/wikinote/internal/server/graceful"
	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/bluemir/wikinote/internal/server/middleware/errors"
)

type Config struct {
	Bind      string
	FrontPage string
}

func NewConfig() *Config {
	return &Config{}
}

type Server struct {
	frontPage string
}

func Run(ctx context.Context, b *backend.Backend, conf *Config) error {
	server := &Server{
		frontPage: conf.FrontPage,
	}

	app := gin.New()

	app.Use(injector.Inject(b))

	// Log
	writer := logrus.WithFields(logrus.Fields{}).WriterLevel(logrus.DebugLevel)
	defer writer.Close()
	app.Use(gin.LoggerWithWriter(writer))

	// Error Handling
	app.Use(errors.Middleware)

	// Recovery
	app.Use(gin.Recovery())
	//app.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
	//	if err, ok := recovered.(string); ok {
	//		c.Error(errors.New(err))
	//		c.Abort()
	//		return
	//	}
	//	c.AbortWithStatus(http.StatusInternalServerError)
	//}))

	app.Use(location.Default(), fixURL)

	// add template
	if html, err := NewRenderer(); err != nil {
		return err
	} else {
		app.SetHTMLTemplate(html)
	}

	// favicon
	app.GET("/favicon.ico", NotFound)

	// Register Routing
	server.route(app, app.NoRoute, b.Plugin)

	logrus.Infof("Run Server on %s", conf.Bind)

	return graceful.Run(ctx, &http.Server{
		Addr:    conf.Bind,
		Handler: app,
	})
}
func NotFound(c *gin.Context) {
	c.String(http.StatusNotFound, "Not Found")
}
