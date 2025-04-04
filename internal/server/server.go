package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-contrib/location"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/backend"
	"github.com/bluemir/wikinote/internal/server/graceful"
	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/bluemir/wikinote/internal/server/middleware/cache"
	"github.com/bluemir/wikinote/internal/server/middleware/errors"
	"github.com/bluemir/wikinote/internal/server/middleware/prom"
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

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()

	app.Use(injector.Inject(b))

	// Log
	writer := logrus.WithFields(logrus.Fields{}).WriterLevel(logrus.InfoLevel)
	defer writer.Close()
	app.Use(gin.LoggerWithWriter(writer))

	// Error Handling
	app.Use(errors.Middleware)

	// Recovery
	app.Use(gin.Recovery())
	app.Use(location.Default(), fixURL)

	// add template
	if html, err := NewHtmlRenderer(); err != nil {
		return err
	} else {
		app.SetHTMLTemplate(html)
	}

	store := cookie.NewStore(xid.New().Bytes())
	store.Options(sessions.Options{
		Path: "/",
	})
	app.Use(sessions.Sessions("wikinote_session", store))

	app.Use(location.Default(), fixURL)
	app.Use(cache.CacheBusting)

	app.Use(prom.Metrics())

	// Register Routing
	server.route(app, app.NoRoute, b.Plugin)

	return graceful.Run(ctx, &http.Server{
		Addr:              conf.Bind,
		Handler:           app,
		ReadHeaderTimeout: 1 * time.Minute,
		WriteTimeout:      3 * time.Minute,
		IdleTimeout:       3 * time.Minute,
	})
}
