package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/injector"
)

func Register(c *gin.Context) error {
	backend := injector.Backends(c)

	req := &struct {
		Username string `form:"username"  validate:"required,min=4"`
		Password string `form:"password"  validate:"required,min=4"`
		Email    string `form:"email"`
		Confirm  string `form:"confirm"   validate:"required,eqfield=Password"`
	}{}
	if err := c.ShouldBind(req); err != nil {
		return err
	}
	if req.Password != req.Confirm {
		return ErrPasswordMismatch
	}

	user, err := backend.Auth.CreateUser(c.Request.Context(), &auth.User{
		Name: req.Username,
		Labels: auth.Labels{
			"wikinote.io/email": req.Email,
		},
	})
	if err != nil {
		return err
	}
	if _, err := backend.Auth.IssueToken(c.Request.Context(), req.Username, req.Password); err != nil {
		return err
	}

	if err := saveUserToSession(c, user); err != nil {
		return err
	}

	c.Redirect(http.StatusSeeOther, "/-/welcome?return="+c.Query("return")) // must use GET method
	return nil
}

func Profile(c *gin.Context) error {
	user, err := User(c)
	if err != nil {
		return err
	}

	c.HTML(http.StatusOK, "system/user/register.html", With(c, struct {
		User *auth.User
	}{
		User: user,
	}))
	return nil
}

func Me(c *gin.Context) error {
	//backend := injector.Backend(c)
	user, err := User(c)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, user)
	return nil
}
