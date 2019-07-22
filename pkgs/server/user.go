package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/pkgs/renderer"
)

func (server *Server) HandleUserList(c *gin.Context) {
	users, err := server.Auth().ListUser()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "/users.html", renderer.Data{"users": users}.With(c))
}

func (server *Server) HandleAPIUserUpdateRole(c *gin.Context) {
	role, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "cannot read body"})
		c.Abort()
		return
	}

	if err := server.Auth().SetUserAttr(c.Param("name"), fmt.Sprintf("rbac/role-%s", role), "true"); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, role)
}
func (server *Server) HandleAPIPutUserAttr(c *gin.Context) {
	key := strings.TrimLeft(c.Param("key"), "/")
	value, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "cannot read body"})
		c.Abort()
		return
	}
	if err := server.Auth().SetUserAttr(c.Param("name"), key, string(value)); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, value)
}
