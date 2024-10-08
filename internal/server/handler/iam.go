package handler

import (
	"net/http"
	"strings"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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
func CreateGroup(c *gin.Context) {
	backend := injector.Backends(c)

	req := struct {
		Name string `form:"name"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if err := backend.Auth.CreateGroup(c.Request.Context(), req.Name); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Redirect(http.StatusSeeOther, "/-/admin/iam/groups")
}
func ListGroups(c *gin.Context) {
	backend := injector.Backends(c)

	groups, err := backend.Auth.ListGroups(c.Request.Context())
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.HTML(http.StatusOK, "admin/iam/groups.html", With(c, KeyValues{
		"groups": groups,
	}))
}
func GetGroup(c *gin.Context) {
	backend := injector.Backends(c)

	group, err := backend.Auth.GetGroup(c.Request.Context(), c.Param("groupName"))
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
func DeleteGroup(c *gin.Context) {
	backend := injector.Backends(c)

	if err := backend.Auth.DeleteGroup(c.Request.Context(), c.Param("groupName")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Redirect(http.StatusSeeOther, "/-/admin/iam/groups")
}
func CreateRole(c *gin.Context) {
	backend := injector.Backends(c)

	req := struct {
		Name string `form:"name"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	role, err := backend.Auth.CreateRole(c.Request.Context(), req.Name, []auth.Rule{})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	logrus.Tracef("%+v", role)

	c.Redirect(http.StatusSeeOther, "/-/admin/iam/roles")

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
func GetRole(c *gin.Context) {
	backend := injector.Backends(c)

	role, err := backend.Auth.GetRole(c.Request.Context(), c.Param("roleName"))
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "admin/iam/role.html", With(c, KeyValues{
		"role": role,
	}))
}
func UpdateRole(c *gin.Context) {
	backend := injector.Backends(c)

	req := struct {
		Rules string `form:"rules"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	rules := []auth.Rule{}
	if err := yaml.Unmarshal([]byte(req.Rules), &rules); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	logrus.Tracef("%+v", rules)

	if err := backend.Auth.UpdateRole(c.Request.Context(), &auth.Role{
		Name:  c.Param("roleName"),
		Rules: rules,
	}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.Redirect(http.StatusSeeOther, "/-/admin/iam/roles/"+c.Param("roleName"))
}
func DeleteRole(c *gin.Context) {
	backend := injector.Backends(c)

	if err := backend.Auth.DeleteRole(c.Request.Context(), c.Param("roleName")); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.Redirect(http.StatusSeeOther, "/-/admin/iam/roles")
}

func CreateAssign(c *gin.Context) {
	backend := injector.Backends(c)

	req := struct {
		Kind string `form:"kind"`
		Name string `form:"name"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if err := backend.Auth.AssignRole(c.Request.Context(), auth.Subject{
		Kind: auth.Kind(req.Kind),
		Name: req.Name,
	}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
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
func GetAssign(c *gin.Context) {
	backend := injector.Backends(c)

	assign, err := backend.Auth.GetAssign(c.Request.Context(), auth.Subject{
		Kind: auth.Kind(c.Param("subjectKind")),
		Name: c.Param("subjectName"),
	})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.HTML(http.StatusOK, "admin/iam/assign.html", With(c, KeyValues{
		"Assign": assign,
	}))
}
func UpdateAssign(c *gin.Context) {
	backend := injector.Backends(c)

	req := struct {
		Roles string `form:"roles"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	assign, err := backend.Auth.GetAssign(c.Request.Context(), auth.Subject{
		Kind: auth.Kind(c.Param("subjectKind")),
		Name: c.Param("subjectName"),
	})
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	assign.Roles = auth.SetFromArray(strings.Split(req.Roles, ","))

	if err := backend.Auth.UpdateAssign(c.Request.Context(), assign); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
}
func DeleteAssign(c *gin.Context) {

	backend := injector.Backends(c)

	if err := backend.Auth.DeleteAssign(c.Request.Context(), auth.Subject{
		Kind: auth.Kind(c.Param("subjectKind")),
		Name: c.Param("subjectName"),
	}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Redirect(http.StatusOK, "/-/admin/iam/assigns")
}
