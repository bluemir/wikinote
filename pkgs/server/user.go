package server

import (
	"net/http"

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
	u, ok, err := server.Auth().GetUser(c.Param("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "user not found"})
		c.Abort()
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		c.Abort()
		return
	}

	/*
		role, err := ioutil.ReadAll(c.Request.Body)
		//u.Role = auth.Role(role)

			if err := Backend(c).Auth().UpdateUser(u); err != nil {
				c.JSON(http.StatusNotFound, gin.H{"msg": err.Error()})
				c.Abort()
				return
			}
	*/
	c.JSON(http.StatusOK, u)
}
