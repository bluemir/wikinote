package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/injector"
)

func (handler *Handler) Login(c *gin.Context) {
	u, err := User(c)
	if errors.Is(err, auth.ErrUnauthorized) {
		c.Error(err)
		c.Abort()
		return
	}

	if u != nil {
		if c.Query("exclude") == "" {
			// logined, but first try.
			c.Redirect(http.StatusTemporaryRedirect, "/-/auth/login?exclude="+u.Name)
			return
		}

		if u.Name == c.Query("exclude") {
			// logined, but try to login same id
			c.Error(auth.ErrUnauthorized)
			c.Abort()
			return
		}
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func Register(c *gin.Context) {
	backend := injector.Backend(c)

	req := &struct {
		Username string `form:"username"  validate:"required,min=4"`
		Password string `form:"password"  validate:"required,min=4"`
		Email    string `form:"email"`
		Confirm  string `form:"confirm"   validate:"required,eqfield=Password"`
	}{}
	if err := c.ShouldBind(req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err := backend.Auth.CreateUser(&auth.User{
		Name: req.Username,
		Labels: auth.Labels{
			"wikinote.io/email": req.Email,
		},
	}); err != nil {
		c.Error(err)
		return
	}
	if _, err := backend.Auth.IssueToken(req.Username, req.Password); err != nil {
		c.Error(err)
		return
	}
	c.Status(http.StatusOK)
}

func Profile(c *gin.Context) {
	user, err := User(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, PageProfile, struct {
		User *auth.User
	}{
		User: user,
	})
}

func (handler *Handler) Me(c *gin.Context) {
	user, err := User(c)
	if errors.Is(err, auth.ErrUnauthorized) {
		c.Error(err)
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, user)
}

func ListUsers(c *gin.Context) {
	backend := injector.Backend(c)

	users, err := backend.Auth.ListUsers()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ListResponse[auth.User]{Items: users})
}
func ListGroups(c *gin.Context) {
	backend := injector.Backend(c)

	groups, err := backend.Auth.ListGroups()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ListResponse[auth.Group]{Items: groups})
}
func ListRoles(c *gin.Context) {
	backend := injector.Backend(c)

	roles, err := backend.Auth.ListRoles()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, ListResponse[auth.Role]{Items: roles})
}
