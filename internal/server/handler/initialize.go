package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/injector"
)

var initCode = xid.New().String()

func RequestInitialize(c *gin.Context) error {
	logrus.Infof("To initailze server, visit %s", c.Request.URL.String()+"/"+initCode)

	c.HTML(http.StatusAccepted, "system/initialize/request.html", With(c, KeyValues{}))
	return nil
}

func Initialze(c *gin.Context) error {
	if c.Param("code") != initCode {
		return ErrInvalidCode
	}
	c.HTML(http.StatusOK, "system/initialize/notice.html", With(c, KeyValues{}))
	return nil
}
func InitialzeAccept(c *gin.Context) error {
	if c.Param("code") != initCode {
		return ErrInvalidCode
	}

	initCode = xid.New().String()

	req := struct {
		Username string
	}{}

	if err := c.ShouldBind(&req); err != nil {
		return err
	}

	backends := injector.Backends(c)

	// Add admin to requested user
	user, _, err := backends.Auth.GetUser(c.Request.Context(), req.Username)
	if err != nil {
		return err
	}

	if user == nil {
		return HttpError{code: http.StatusBadRequest, message: "user not found"}
	}
	user.AddGroup("admin")

	if err := backends.Auth.UpdateUser(c.Request.Context(), user); err != nil {
		return err
	}

	// reset admin role
	if err := backends.Auth.UpdateRole(c.Request.Context(), &auth.Role{
		Name: "admin",
		Rules: []auth.Rule{
			{}, // allow all
		},
	}); err != nil {
		return err
	}

	c.HTML(http.StatusOK, "system/initialize/done.html", With(c, nil))
	return nil
}
