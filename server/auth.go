package server

import (
	"net/http"

	"github.com/bluemir/go-utils/auth"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/server/renderer"
)

func BasicAuth(c *gin.Context) {
	token, err := Backend(c).Auth().HttpAuth(c.GetHeader("Authorization"))

	switch err.Code() {
	case auth.ErrorNone:
		logrus.Debug("Login user :", token.Username)
		c.SetCookie("logined", token.Username, 0, "", "", false, true)
		c.Set(TOKEN, token)
		return
	case auth.ErrorEmptyHeader:
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
	case auth.ErrorWrongEncoding, auth.ErrorBadToken:
		FlashMessage(c).Warn("Connot decode auth token")
		c.HTML(http.StatusBadRequest, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	case auth.ErrorEmptyAccount:
		logrus.Debug("Empty Account")
		c.SetCookie("logined", "", -1, "", "", false, true)
		c.Header("WWW-Authenticate", AuthenicateString)
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return // just pass, it is a guest or logout user
	case auth.ErrorUnauthorized:
		logrus.Debug("unauthorized")
		FlashMessage(c).Warn("Error on auth, id password not matched")
		c.Header("WWW-Authenticate", AuthenicateString)
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	default:
		FlashMessage(c).Warn("Somethings Wrong. plz contact system admin")
		c.HTML(http.StatusInternalServerError, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}
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
	if registeForm.Password != registeForm.Confirm {
		FlashMessage(c).Warn("password not confirm...")
		c.Redirect(http.StatusSeeOther, "/!/auth/register")
		return
	}

	err = Backend(c).Auth().CreateUser(&auth.User{
		Name: registeForm.Id,
		Labels: map[string]string{
			"email": registeForm.Email,
		},
	})
	if err != nil {
		FlashMessage(c).Warn("fail to register: %s", err.Error())
		c.Redirect(http.StatusSeeOther, "/!/auth/register")
		return
	}
	_, err = Backend(c).Auth().IssueToken(registeForm.Id, registeForm.Password)
	if err != nil {
		FlashMessage(c).Warn("fail to register: %s", err.Error())
		c.Redirect(http.StatusSeeOther, "/!/auth/register")
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func HandleLogin(c *gin.Context) {
	//Backend(c).Auth().HttpAuth(c.GetHeader("Authorization"))
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
		c.SetCookie("logined", "", -1, "", "", false, true)
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
