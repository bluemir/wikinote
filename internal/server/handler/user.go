package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

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
func (handler *Handler) Can(c *gin.Context) {
	user, err := User(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	verb := c.Param("verb")
	kind := strings.TrimPrefix(c.Param("kind"), "/")

	logrus.WithField("verb", verb).WithField("kind", kind).Trace("API called")

	if err := handler.backend.Auth.Can(user, auth.Verb(verb), auth.KeyValues{"kind": kind}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusOK)
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
