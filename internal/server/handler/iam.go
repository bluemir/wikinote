package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	backend := injector.Backends(c)

	users, err := backend.Auth.ListUsers(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "admin/iam/users.html", with(c, KeyValues{
		"users": users,
	}))
}
func ListGroups(c *gin.Context) {
	backend := injector.Backends(c)

	groups, err := backend.Auth.ListGroups(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "admin/iam/groups.html", with(c, KeyValues{
		"groups": groups,
	}))
}
func ListRoles(c *gin.Context) {
	backend := injector.Backends(c)

	roles, err := backend.Auth.ListRoles(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "admin/iam/roles.html", with(c, KeyValues{
		"roles": roles,
	}))
}
