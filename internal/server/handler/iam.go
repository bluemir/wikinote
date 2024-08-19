package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	backend := injector.Backend(c)

	users, err := backend.Auth.ListUsers(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "admin/iam/users.html", renderData(c, KeyValues{
		"users": users,
	}))
}
func ListGroups(c *gin.Context) {
	backend := injector.Backend(c)

	groups, err := backend.Auth.ListGroups(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "admin/iam/groups.html", renderData(c, KeyValues{
		"groups": groups,
	}))
}
func ListRoles(c *gin.Context) {
	backend := injector.Backend(c)

	roles, err := backend.Auth.ListRoles(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "admin/iam/roles.html", renderData(c, KeyValues{
		"roles": roles,
	}))
}
