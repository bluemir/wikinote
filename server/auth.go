package server

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend/user"
	"github.com/bluemir/wikinote/server/renderer"
)

func BasicAuth(c *gin.Context) {
	str := c.GetHeader("Authorization")
	if str == "" {
		_, err := c.Cookie("logined")
		if err != http.ErrNoCookie {
			logrus.Debug("cookie found but auth not found")
			c.Header("WWW-Authenticate", AuthenicateString)
			c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
			c.Abort()
			return
		}
		// Skip auth
		logrus.Debug("skip auth")
		return
	}

	arr := strings.SplitN(str, " ", 2)

	if len(arr) != 2 {
		c.HTML(http.StatusBadRequest, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}

	if arr[0] != "Basic" {
		FlashMessage(c).Warn("Token type is not 'Basic'")
		c.HTML(http.StatusBadRequest, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}
	buf, err := base64.StdEncoding.DecodeString(arr[1])
	if err != nil {
		FlashMessage(c).Warn("Connot decode auth token")
		c.HTML(http.StatusBadRequest, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}

	arr = strings.SplitN(string(buf), ":", 2)
	username := arr[0]
	password := ""
	if len(arr) > 1 {
		password = arr[1]
	}
	if username == "" && password == "" {
		c.SetCookie("logined", "", -1, "", "", false, true)
		return // just pass, it is a guest
	}
	user, ok, err := Backend(c).User().Auth(username, password)
	if err != nil {
		FlashMessage(c).Warn("Somethings Wrong. plz contact system admin")
		c.HTML(http.StatusInternalServerError, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}
	if !ok {
		FlashMessage(c).Warn("Error on auth, id password not matched")
		c.Header("WWW-Authenticate", AuthenicateString)
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}
	logrus.Debug("Login user :", user)
	c.SetCookie("logined", user.Name, 0, "", "", false, true)
	c.Set(USER, user)
	return
}
func HandleRegisterForm(c *gin.Context) {
	c.HTML(http.StatusOK, "/register.html", renderer.Data{}.With(c))
}
func HandleRegister(c *gin.Context) {
	registeForm := &struct {
		Id       string `form:"id"`
		Password string `form:"password"`
		Email    string `form:"email"`
		Confirm  string `form:"confirm"`
	}{}
	err := c.ShouldBind(registeForm)
	if err != nil {
		FlashMessage(c).Warn("bad Request")
		c.Redirect(http.StatusSeeOther, "/!/auth/register")
		c.Abort()
		return
	}

	logrus.Debugf("%+v", registeForm)
	if registeForm.Password != registeForm.Confirm {
		FlashMessage(c).Warn("password not confirm...")
		c.Redirect(http.StatusSeeOther, "/!/auth/register")
		return
	}
	err = Backend(c).User().New(&user.User{
		Name:  registeForm.Id,
		Email: registeForm.Email,
	}, registeForm.Password)
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
		c.Header("WWW-Authenticate", AuthenicateString)
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
func HandleLogout(c *gin.Context) {
	str := c.GetHeader("Authorization")
	if str != "Basic Og==" { // empty id pass word
		c.Header("WWW-Authenticate", AuthenicateString)
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
