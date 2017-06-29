package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend/user"
)

func HandleRegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "/register.html", Data(c))
}
func HandleRegister(c *gin.Context) {
	registeForm := &struct {
		Id       string `form:"id"`
		Password string `form:"password"`
		Email    string `form:"email"`
		Confirm  string `form:"confirm"`
	}{}
	err := c.Bind(registeForm)
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Debugf("%+v", registeForm)
	if registeForm.Password != registeForm.Confirm {
		FlashMessage(c).Warn("password not confirm...")
		c.Redirect(http.StatusSeeOther, "/!/auth/register")
		return
	}
	err = Backend(c).User().New(&user.User{
		Id:       registeForm.Id,
		Email:    registeForm.Email,
		Password: user.NewPassword(registeForm.Password),
	})
	if err != nil {
		FlashMessage(c).Warn("fail to register: %s", err.Error())
		c.Redirect(http.StatusSeeOther, "/!/auth/register")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}
