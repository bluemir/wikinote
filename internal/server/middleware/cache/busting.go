package cache

import (
	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/buildinfo"
)

const (
	REVVED = "__REVVED__"
)

var (
	rev = buildinfo.Signature()[:16]
)

func CacheBusting(c *gin.Context) {
	c.Set(REVVED, rev)
}

func Rev() string {
	return rev
}
