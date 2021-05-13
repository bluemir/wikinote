package server

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) HandleAPIUserUpdateRole(c *gin.Context) {
	role, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "cannot read body"})
		c.Abort()
		return
	}

	username := c.Param("name")
	u, ok, err := server.Auth.GetUser(username)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"msg": "not found"})
		c.Abort()
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}
	if u.Labels == nil {
		u.Labels = map[string]string{}
	}
	u.Labels["role/"+string(role)] = "true"

	if err := server.Auth.UpdateUser(u); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, u)
}
