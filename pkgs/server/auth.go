package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkgs/auth"
	"github.com/bluemir/wikinote/pkgs/backend"
	"github.com/bluemir/wikinote/pkgs/renderer"
)

func isLogined(c *gin.Context) bool {
	_, err := c.Cookie("logined")
	if err == http.ErrNoCookie {
		return false
	}
	return true
}

func Token(c *gin.Context) *auth.Token {
	token, ok := c.Get(TOKEN)
	if ok {
		return token.(*auth.Token)
	}
	return nil
}

func BasicAuthn(c *gin.Context) {
	token, err := Backend(c).Auth().HttpAuth(c.GetHeader("Authorization"))

	switch auth.ErrorCode(err) {
	case auth.ErrNone:
		logrus.Debug("Login user :", token.UserName)
		c.SetCookie("logined", token.UserName, 0, "", "", false, true)
		c.Set(TOKEN, token)
		return
	case auth.ErrEmptyAccount: // it means logout
		logrus.Debug("Empty Account")
		if isLogined(c) {
			c.SetCookie("logined", "", -1, "", "", false, true)
		} else {
			c.Header("WWW-Authenticate", AuthenicateString)
			c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
			c.Abort()
		}
		return
	case auth.ErrEmptyHeader:
		logrus.Debug("Empty header")
		if isLogined(c) {
			logrus.Debug("cookie found but auth not found")
			c.Header("WWW-Authenticate", AuthenicateString)
			c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
			c.Abort()
			return
		}
		// Skip auth
		logrus.Debug("skip auth")
		return
	case auth.ErrWrongEncoding, auth.ErrBadToken:
		FlashMessage(c).Warn("Connot decode auth token")
		c.HTML(http.StatusBadRequest, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	case auth.ErrUnauthorized:
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
func Authz(actions ...string) func(c *gin.Context) {
	return func(c *gin.Context) {
		token, ok := c.Get(TOKEN)
		if !ok {
			c.Header("WWW-Authenticate", AuthenicateString)
			c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
			c.Abort()
			return
		}

		subject := Backend(c).Auth().Subject(token.(*auth.Token))
		object := &backend.AuthzObject{Backend(c).File().Attr(c.Request.URL.Path)}

		for _, action := range actions {
			result, err := Backend(c).Plugin().AuthCheck(&auth.Context{
				Subject: subject,
				Object:  object,
				Action:  action,
			})
			if err != nil {
				c.HTML(http.StatusInternalServerError, "errors/internal.html", renderer.Data{}.With(c))
				c.Abort()
				return
			}
			switch result {
			case auth.Reject:
				c.HTML(http.StatusForbidden, "/errors/forbidden.html", renderer.Data{}.With(c))
				c.Abort()
				return
			case auth.Accept:
				// check next
				continue
			case auth.Unknown:
				c.HTML(http.StatusInternalServerError, "errors/internal.html", renderer.Data{}.With(c))
				c.Abort()
				return
			}
		}
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
	})
	//Backend(c).Auth().SetAttr(username, HandleRegisterForm.Email)
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
		c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", renderer.Data{}.With(c))
		c.Abort()
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
