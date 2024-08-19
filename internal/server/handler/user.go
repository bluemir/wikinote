package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/injector"
)

func Login(c *gin.Context) {
	u, err := User(c)
	if errors.Is(err, auth.ErrUnauthorized) {
		c.Error(err)
		c.Abort()
		return
	}

	exclude, exist := c.GetQuery("exclude")
	if u != nil && exist {
		if exclude == "" {
			// logined, try to login other id.
			c.Redirect(http.StatusTemporaryRedirect, "/-/auth/login?exclude="+u.Name)
			return
		}
		if u.Name == exclude {
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
	user, err := backend.Auth.CreateUser(c.Request.Context(), &auth.User{
		Name: req.Username,
		Labels: auth.Labels{
			"wikinote.io/email": req.Email,
		},
	})
	if err != nil {
		c.Error(err)
		return
	}
	if _, err := backend.Auth.IssueToken(c.Request.Context(), req.Username, req.Password); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
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

func Me(c *gin.Context) {
	//backend := injector.Backend(c)
	user, err := User(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, user)
}
