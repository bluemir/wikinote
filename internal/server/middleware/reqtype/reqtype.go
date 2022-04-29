package reqtype

import (
	"mime"

	"github.com/gin-gonic/gin"
)

type ReqType int

const (
	Unknown = iota
	API
	HTML
)

func FindRequestType(c *gin.Context) ReqType {
	for _, ct := range c.Accepted {
		mt, _, err := mime.ParseMediaType(ct)
		if err != nil {
			continue
		}
		switch mt {
		case "application/json":
			return API
		case "text/html":
			return HTML
		}
	}
	return Unknown
}

func MarkAPI(c *gin.Context) {
	c.SetAccepted("application/json")
}
func MakrHTML(c *gin.Context) {
	c.SetAccepted("text/html")
}
