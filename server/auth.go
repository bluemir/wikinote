package server

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend/user"
)

func BasicAuth(c *gin.Context) {
	str := c.GetHeader("Authorization")
	if str == "" {
		//c.Header("WWW-Authenticate", "Basic realm=\"Auth required!\"")
		//c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", Data(c))
		//c.Abort()
		return
	}

	if str[:len("Basic ")] != "Basic " {
		c.HTML(http.StatusBadRequest, "/errors/unauthorized.html", Data(c))
		c.Abort()
		return
	}
	buf, err := base64.StdEncoding.DecodeString(str[len("Basic "):])
	if err != nil {
		c.HTML(http.StatusBadRequest, "/errors/unauthorized.html", Data(c))
		c.Abort()
		return
	}
	arr := strings.SplitN(string(buf), ":", 2)
	username := arr[0]
	password := ""
	if len(arr) > 1 {
		password = arr[1]
	}

	user, err := Backend(c).User().Get(username)
	if err != nil {
		FlashMessage(c).Warn("Error on get user: %s", err.Error())
		c.Header("WWW-Authenticate", "Basic realm=\"Auth required!\"")
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", Data(c))
		c.Abort()
		return
	}

	if user == nil {
		FlashMessage(c).Warn("Error on get user: %s", err.Error())
		c.Header("WWW-Authenticate", "Basic realm=\"Auth required!\"")
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", Data(c))
		c.Abort()
		return

	}
	if !user.Password.Check(password) {
		FlashMessage(c).Warn("wrong password")
		c.Header("WWW-Authenticate", "Basic realm=\"Auth required!\"")
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", Data(c))
		c.Abort()
		return
	}
	c.Set(USERNAME, username)
	return

}
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

func HandleLogin(c *gin.Context) {
	str := c.GetHeader("Authorization")
	if str == "" {
		c.Header("WWW-Authenticate", "Basic realm=\"Auth required!\"")
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", Data(c))
		c.Abort()
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, Backend(c).Config().FrontPage)
}
