package server

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkg/backend"
)

type Config struct {
	Bind       string
	TLSDomains []string
}

func Run(b *backend.Backend, conf *Config) error {
	server := &Server{b}

	app := gin.New()

	// Log
	writer := logrus.WithField("from", "gin").Writer()
	defer writer.Close()
	app.Use(gin.LoggerWithWriter(writer))

	// Recovery
	app.Use(gin.Recovery())

	// Session
	store := cookie.NewStore([]byte("__wikinote__"))
	app.Use(sessions.Sessions("session", store))

	// Renderer
	if html, err := NewRenderer(); err != nil {
		return errors.WithStack(err)
	} else {
		app.SetHTMLTemplate(html)
	}

	server.RegisterRoute(app)

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
