package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/auth"
)

func (server *Server) HandleRegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "/register.html", gin.H{})
}
func (server *Server) HandleRegister(c *gin.Context) {
	req := &struct {
		Name     string `form:"name"`
		Password string `form:"password"`
		Email    string `form:"email"`
		Confirm  string `form:"confirm"`
	}{}

	if err := c.ShouldBind(req); err != nil {
		// Set cookie?
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

	err := server.Auth.CreateUser(req.Name, auth.Labels{
		"wikinote.io/email": req.Email,
		//"role/default":      "true",
	})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{
			"retryURL": "/!/auth/register",
			"message":  "fail to register new user try again",
		})
		c.Abort()
		return
	}

	_, err = server.Auth.IssueToken(req.Name, req.Password)
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
