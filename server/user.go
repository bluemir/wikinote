package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUserList(c *gin.Context) {
	users, err := Backend(c).User().List()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "users.html", Data(c).Set("users", users))
}
