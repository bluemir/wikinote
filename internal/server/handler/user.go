package handler

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/injector"
)

func Register(c *gin.Context) {
	backend := injector.Backends(c)

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
	if req.Password != req.Confirm {
		c.AbortWithError(http.StatusBadRequest, errors.New("password & confirm password note matched"))
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

	session := sessions.Default(c)
	session.Set(SessionKeyUser, user)

	if err := session.Save(); err != nil {
		c.Error(err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/-/welcome?return="+c.Query("return")) // must use GET method
}

func Profile(c *gin.Context) {
	user, err := User(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, PageProfile, With(c, struct {
		User *auth.User
	}{
		User: user,
	}))
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
