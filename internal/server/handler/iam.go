package handler

import (
	"net/http"
	"strings"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func GetUser(c *gin.Context) {
	backend := injector.Backends(c)

	user, _, err := backend.Auth.GetUser(c.Request.Context(), c.Param("username"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "admin/iam/user.html", With(c, KeyValues{
		"user": user,
	}))

}
func ListUsers(c *gin.Context) {
	backend := injector.Backends(c)

	users, err := backend.Auth.ListUsers(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}

	c.HTML(http.StatusOK, "admin/iam/users.html", With(c, KeyValues{
		"users": users,
	}))
}
func UpdateUser(c *gin.Context) {
	backend := injector.Backends(c)

	req := struct {
		Groups string
	}{}

	user, _, err := backend.Auth.GetUser(c.Request.Context(), c.Param("username"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	for _, group := range strings.Split(req.Groups, ",") {
		user.AddGroup(group)
	}

	if err := backend.Auth.UpdateUser(c.Request.Context(), user); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
}
func ListGroups(c *gin.Context) {
	backend := injector.Backends(c)

	groups, err := backend.Auth.ListGroups(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "admin/iam/groups.html", With(c, KeyValues{
		"groups": groups,
	}))
}
func GetGroup(c *gin.Context) {
	backend := injector.Backends(c)

	group, err := backend.Auth.GetGroup(c.Request.Context(), c.Param("groupname"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	logrus.Tracef("%+v", group)
	c.HTML(http.StatusOK, "admin/iam/group.html", With(c, KeyValues{
		"group": group,
	}))
}
func ListRoles(c *gin.Context) {
	backend := injector.Backends(c)

	roles, err := backend.Auth.ListRoles(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "admin/iam/roles.html", With(c, KeyValues{
		"roles": roles,
	}))
}

func ListAssigns(c *gin.Context) {
	backend := injector.Backends(c)

	assigns, err := backend.Auth.ListAssigns(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "admin/iam/assigns.html", With(c, KeyValues{
		"assigns": assigns,
	}))
}
