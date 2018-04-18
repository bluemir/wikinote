package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/server/renderer"
)

func HandleUserList(c *gin.Context) {
	users, err := Backend(c).User().List()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "users.html", renderer.Data{"users": users}.With(c))
}
