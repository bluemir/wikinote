package server

import (
	"context"
	"net/http"

	"github.com/gin-contrib/location"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
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

	// Register Routing
	server.route(app, app.NoRoute, b.Plugin)

	return graceful.Run(ctx, &http.Server{
		Addr:    conf.Bind,
		Handler: app,
	})
}
