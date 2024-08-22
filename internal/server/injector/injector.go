package injector

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"github.com/bluemir/wikinote/internal/backend"
)

var keyBackend = xid.New().String()

func Inject(b *backend.Backend) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(keyBackend, b)
	}
}
func Backends(c *gin.Context) *backend.Backend {
	return c.MustGet(keyBackend).(*backend.Backend)
}
