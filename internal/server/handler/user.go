package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/gin-gonic/gin"
)

func (handler *Handler) RegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "/register.html", gin.H{})
}

func (handler *Handler) Register(c *gin.Context) {
	req := &struct {
		Name     string `form:"name"`
		Password string `form:"password"`
		Email    string `form:"email"`
		Confirm  string `form:"confirm"`
	}{}

	if err := c.ShouldBind(req); err != nil {
		c.HTML(http.StatusBadRequest, "/errors/bad-request.html", gin.H{
			"retryURL": "/!/auth/register",
		})
		c.Abort()
		return
	}
	if req.Password != req.Confirm {
		c.HTML(http.StatusBadRequest, "/errors/bad-request.html", gin.H{
			"retryURL": "/!/auth/register",
			"message":  "password and password confirm not matched",
		})
		c.Abort()
		return
	}

	err := handler.backend.Auth.CreateUser(req.Name, auth.Labels{
		"wikinote.io/email": req.Email,
	})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{
			"retryURL": "/!/auth/register",
			"message":  "fail to register new user try again",
		})
		c.Abort()
		return
	}

	_, err = handler.backend.Auth.IssueToken(req.Name, req.Password, nil)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "errors/internal-server-error.html", gin.H{
			"retryURL": "/!/auth/register",
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
		HTTPErrorHandler(c, err, WithAuthHeader)
		return
	}

	c.HTML(http.StatusOK, "/profile.html", gin.H{
		"user": user,
	})
}
