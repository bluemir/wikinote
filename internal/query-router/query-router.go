package queryrouter

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func New() QueryRouter {
	return &router{
		handlers: map[string]map[string][]gin.HandlerFunc{

			http.MethodGet:    map[string][]gin.HandlerFunc{},
			http.MethodPost:   map[string][]gin.HandlerFunc{},
			http.MethodPut:    map[string][]gin.HandlerFunc{},
			http.MethodDelete: map[string][]gin.HandlerFunc{},
		},
	}
}

type QueryRouter interface {
	GET(query string, handlers ...gin.HandlerFunc)
	POST(query string, handlers ...gin.HandlerFunc)
	PUT(query string, handlers ...gin.HandlerFunc)
	DELETE(query string, handlers ...gin.HandlerFunc)

	Register(method string, query string, handlers ...gin.HandlerFunc)
	Handler(c *gin.Context)
}

type router struct {
	handlers map[string]map[string][]gin.HandlerFunc
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
func (r *router) GET(query string, handlers ...gin.HandlerFunc) {
	r.Register(http.MethodGet, query, handlers...)
}
func (r *router) POST(query string, handlers ...gin.HandlerFunc) {
	r.Register(http.MethodPost, query, handlers...)
}
func (r *router) PUT(query string, handlers ...gin.HandlerFunc) {
	r.Register(http.MethodPut, query, handlers...)
}
func (r *router) DELETE(query string, handlers ...gin.HandlerFunc) {
	r.Register(http.MethodDelete, query, handlers...)
}
func (r *router) Register(method string, query string, handlers ...gin.HandlerFunc) {
	// TODO dup check
	// TODO nil check
	r.handlers[method][query] = handlers
}
