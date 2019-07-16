package queryrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func New() QueryRouter {
	return &router{
		handlers: map[string]map[string][]func(c *gin.Context){
			http.MethodGet:    map[string][]func(c *gin.Context){},
			http.MethodPost:   map[string][]func(c *gin.Context){},
			http.MethodPut:    map[string][]func(c *gin.Context){},
			http.MethodDelete: map[string][]func(c *gin.Context){},
		},
	}
}

type QueryRouter interface {
	Register(method string, query string, handlers ...func(c *gin.Context))
	Handler(c *gin.Context)
}

type router struct {
	handlers map[string]map[string][]func(c *gin.Context)
}

func (r *router) Handler(c *gin.Context) {
	log := logrus.
		WithField("method", "queryrotuer.Handler").
		WithField("method", c.Request.Method)

	log.Trace("start route")

	qh, ok := r.handlers[c.Request.Method]
	if !ok {
		log.Debug("method not found")
		return
	}

	for k, hs := range qh {
		if k == "*" {
			continue // '*' will handle at last
		}
		if _, ok := c.GetQuery(k); ok {
			for _, h := range hs {
				if c.IsAborted() {
					return
				}
				h(c)
			}
			return
		}
	}
	if hs, ok := qh["*"]; ok {
		for _, h := range hs {
			if c.IsAborted() {
				return
			}
			h(c)
		}
	}
}
func (r *router) Register(method string, query string, handlers ...func(c *gin.Context)) {
	// TODO dup check
	// TODO nil check
	r.handlers[method][query] = handlers
}
