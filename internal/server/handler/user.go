package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/auth"
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

func (handler *Handler) Register(c *gin.Context) {
	req := &struct {
		Username string `form:"username"`
		Password string `form:"password"`
		Email    string `form:"email"`
		Confirm  string `form:"confirm"`
	}{}

	if err := c.ShouldBind(req); err != nil {
		c.HTML(http.StatusBadRequest, "/errors/bad-request.html", gin.H{
			"retryURL": "/-/auth/register",
		})
		return
	}
	if req.Password != req.Confirm {
		c.HTML(http.StatusBadRequest, "/errors/bad-request.html", gin.H{
			"retryURL": "/-/auth/register",
			"message":  "password and password confirm not matched",
		})
		return
	}

	if err := handler.backend.Auth.CreateUser(&auth.User{
		Name: req.Username,
		Labels: auth.Labels{
			"wikinote.io/email": req.Email,
		},
	}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if _, err := handler.backend.Auth.IssueToken(req.Username, req.Password, nil); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "/welcome.html", gin.H{})
}
func (handler *Handler) Profile(c *gin.Context) {
	user, err := User(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "/profile.html", gin.H{
		"user": user,
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
	kind := c.Param("kind")

	if err := handler.backend.Auth.Can(user, auth.Verb(verb), auth.KeyValues{"kind": kind}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{})
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
