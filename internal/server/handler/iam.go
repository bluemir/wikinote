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

func GetUser(c *gin.Context) error {
	backend := injector.Backends(c)

	user, _, err := backend.Auth.GetUser(c.Request.Context(), c.Param("username"))
	if err != nil {
		return err
	}

	c.HTML(http.StatusOK, "admin/iam/user.html", With(c, KeyValues{
		"user": user,
	}))
	return nil
}
func ListUsers(c *gin.Context) error {
	backend := injector.Backends(c)

	users, err := backend.Auth.ListUsers(c.Request.Context())
	if err != nil {
		return err
	}

	c.HTML(http.StatusOK, "admin/iam/users.html", With(c, KeyValues{
		"users": users,
	}))
	return nil
}
func UpdateUser(c *gin.Context) error {
	backend := injector.Backends(c)

	req := struct {
		Groups string
	}{}

	user, _, err := backend.Auth.GetUser(c.Request.Context(), c.Param("username"))
	if err != nil {
		return err
	}

	for group := range strings.SplitSeq(req.Groups, ",") {
		user.AddGroup(group)
	}

	if err := backend.Auth.UpdateUser(c.Request.Context(), user); err != nil {
		return err
	}
	return nil
}
func CreateGroup(c *gin.Context) error {
	backend := injector.Backends(c)

	req := struct {
		Name string `form:"name"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		return err
	}

	if err := backend.Auth.CreateGroup(c.Request.Context(), req.Name); err != nil {
		return err
	}

	c.Redirect(http.StatusSeeOther, "/-/admin/iam/groups")
	return nil
}
func ListGroups(c *gin.Context) error {
	backend := injector.Backends(c)

	groups, err := backend.Auth.ListGroups(c.Request.Context())
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "admin/iam/groups.html", With(c, KeyValues{
		"groups": groups,
	}))
	return nil
}
func GetGroup(c *gin.Context) error {
	backend := injector.Backends(c)

	group, err := backend.Auth.GetGroup(c.Request.Context(), c.Param("groupName"))
	if err != nil {
		return err
	}
	logrus.Tracef("%+v", group)
	c.HTML(http.StatusOK, "admin/iam/group.html", With(c, KeyValues{
		"group": group,
	}))
	return nil
}
func DeleteGroup(c *gin.Context) error {
	backend := injector.Backends(c)

	if err := backend.Auth.DeleteGroup(c.Request.Context(), c.Param("groupName")); err != nil {
		return err
	}

	c.Redirect(http.StatusSeeOther, "/-/admin/iam/groups")
	return nil
}
func CreateRole(c *gin.Context) error {
	backend := injector.Backends(c)

	req := struct {
		Name string `form:"name"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		return err
	}
	role, err := backend.Auth.CreateRole(c.Request.Context(), req.Name, []auth.Rule{})
	if err != nil {
		return err
	}

	logrus.Tracef("%+v", role)

	c.Redirect(http.StatusSeeOther, "/-/admin/iam/roles")
	return nil
}
func ListRoles(c *gin.Context) error {
	backend := injector.Backends(c)

	roles, err := backend.Auth.ListRoles(c.Request.Context())
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "admin/iam/roles.html", With(c, KeyValues{
		"roles": roles,
	}))
	return nil
}
func GetRole(c *gin.Context) error {
	backend := injector.Backends(c)

	role, err := backend.Auth.GetRole(c.Request.Context(), c.Param("roleName"))
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "admin/iam/role.html", With(c, KeyValues{
		"role": role,
	}))
	return nil
}
func UpdateRole(c *gin.Context) error {
	backend := injector.Backends(c)

	req := struct {
		Rules string `form:"rules"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		return err
	}

	rules := []auth.Rule{}
	if err := yaml.Unmarshal([]byte(req.Rules), &rules); err != nil {
		return err
	}
	logrus.Tracef("%+v", rules)

	if err := backend.Auth.UpdateRole(c.Request.Context(), &auth.Role{
		Name:  c.Param("roleName"),
		Rules: rules,
	}); err != nil {
		return err
	}
	c.Redirect(http.StatusSeeOther, "/-/admin/iam/roles/"+c.Param("roleName"))
	return nil
}
func DeleteRole(c *gin.Context) error {
	backend := injector.Backends(c)

	if err := backend.Auth.DeleteRole(c.Request.Context(), c.Param("roleName")); err != nil {
		return err
	}
	c.Redirect(http.StatusSeeOther, "/-/admin/iam/roles")
	return nil
}

func CreateAssign(c *gin.Context) error {
	backend := injector.Backends(c)

	req := struct {
		Kind string `form:"kind"`
		Name string `form:"name"`
	}{}
	if err := c.ShouldBind(&req); err != nil {
		return err
	}

	if err := backend.Auth.AssignRole(c.Request.Context(), auth.Subject{
		Kind: auth.Kind(req.Kind),
		Name: req.Name,
	}); err != nil {
		return err
	}

	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
	return nil
}
func ListAssigns(c *gin.Context) error {
	backend := injector.Backends(c)

	assigns, err := backend.Auth.ListAssigns(c.Request.Context())
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "admin/iam/assigns.html", With(c, KeyValues{
		"assigns": assigns,
	}))
	return nil
}
func GetAssign(c *gin.Context) error {
	backend := injector.Backends(c)

	assign, err := backend.Auth.GetAssign(c.Request.Context(), auth.Subject{
		Kind: auth.Kind(c.Param("subjectKind")),
		Name: c.Param("subjectName"),
	})
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "admin/iam/assign.html", With(c, KeyValues{
		"Assign": assign,
	}))
	return nil
}
func UpdateAssign(c *gin.Context) error {
	backend := injector.Backends(c)

	req := struct {
		Roles string `form:"roles"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		return err
	}

	assign, err := backend.Auth.GetAssign(c.Request.Context(), auth.Subject{
		Kind: auth.Kind(c.Param("subjectKind")),
		Name: c.Param("subjectName"),
	})
	if err != nil {
		return err
	}

	assign.Roles = auth.SetFromArray(strings.Split(req.Roles, ","))

	if err := backend.Auth.UpdateAssign(c.Request.Context(), assign); err != nil {
		return err
	}
	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
	return nil
}
func DeleteAssign(c *gin.Context) error {

	backend := injector.Backends(c)

	if err := backend.Auth.DeleteAssign(c.Request.Context(), auth.Subject{
		Kind: auth.Kind(c.Param("subjectKind")),
		Name: c.Param("subjectName"),
	}); err != nil {
		return err
	}

	c.Redirect(http.StatusOK, "/-/admin/iam/assigns")
	return nil
}
