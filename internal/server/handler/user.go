package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/auth"
)

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
		c.Abort()
		return
	}
	if req.Password != req.Confirm {
		c.HTML(http.StatusBadRequest, "/errors/bad-request.html", gin.H{
			"retryURL": "/-/auth/register",
			"message":  "password and password confirm not matched",
		})
		c.Abort()
		return
	}

	err := handler.backend.Auth.CreateUser(&auth.User{
		Name: req.Username,
		Labels: auth.Labels{
			"wikinote.io/email": req.Email,
		},
	})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{
			"retryURL": "/-/auth/register",
			"message":  "fail to register new user try again",
			"error":    err.Error(),
		})
		logrus.Error(err)
		c.Abort()
		return
	}

	_, err = handler.backend.Auth.IssueToken(req.Username, req.Password, nil)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/internal-server-error.html", gin.H{
			"retryURL": "/-/auth/register",
			"message":  "fail to register new user try again",
		})
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
