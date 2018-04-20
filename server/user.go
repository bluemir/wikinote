package server

import (
	"io/ioutil"
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

func HandleAPIUserUpdateRole(c *gin.Context) {
	u, err := Backend(c).User().Get(c.Param("name"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	role, err := ioutil.ReadAll(c.Request.Body)
	u.Role = string(role)
	if err := Backend(c).User().Put(u); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, u)
}
