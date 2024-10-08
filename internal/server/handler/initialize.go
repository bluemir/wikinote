package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/injector"
)

var code = xid.New().String()

func RequestInitialize(c *gin.Context) {
	logrus.Infof("To initailze server, visit %s", c.Request.URL.String()+"/"+code)

	c.HTML(http.StatusAccepted, "system/initialize/request.html", With(c, KeyValues{}))
}

func Initialze(c *gin.Context) {
	if c.Param("code") != code {
		c.Error(HttpError{code: http.StatusForbidden, message: "code not matched"})
		return
	}
	c.HTML(http.StatusOK, "system/initialize/notice.html", With(c, KeyValues{}))
}
func InitialzeAccept(c *gin.Context) {
	if c.Param("code") != code {
		c.Error(errors.New("code not matched"))
		return
	}

	code = xid.New().String()

	req := struct {
		Username string
	}{}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		return
	}

	backends := injector.Backends(c)

	// Add admin to requested user
	user, _, err := backends.Auth.GetUser(c.Request.Context(), req.Username)
	if err != nil {
		c.Error(err)
		return
	}

	if user == nil {
		c.Error(HttpError{code: http.StatusBadRequest, message: "user not found"})
		//c.Error(errors.New("user not found"))
		return
	}
	user.AddGroup("admin")

	if err := backends.Auth.UpdateUser(c.Request.Context(), user); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	// reset admin role
	if err := backends.Auth.UpdateRole(c.Request.Context(), &auth.Role{
		Name: "admin",
		Rules: []auth.Rule{
			{}, // allow all
		},
	}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "system/initialize/done.html", With(c, nil))
}
