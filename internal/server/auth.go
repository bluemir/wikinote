package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/auth"
)

func (server *Server) Authn(c *gin.Context) {
	user, err := server.Backend.Auth.HTTP(c.Request)

	if err != nil {
		switch err {
		case auth.ErrEmptyHeader:
			// check Authz with guest
			// if failed, must request login
			return
		default:
			c.Header("WWW-Authenticate", AuthenicateString)
			c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", gin.H{})
			c.Abort()
			return
		}
	}
	logrus.Tracef("auth as '%s'", user.Name)

	// check Authz with token

	c.Set(SessionKeyUser, user)

	return // next

}

func getUser(c *gin.Context) (*auth.User, error) {
	u, exist := c.Get(SessionKeyUser)
	if !exist {
		return nil, auth.ErrUnauthorized // TODO
	}
	user, ok := u.(*auth.User)
	if !ok {
		return nil, auth.ErrUnauthorized // TODO
	}
	return user, nil
}

type ResourceGetter func(c *gin.Context) (auth.Resource, error)

func (server *Server) Authz(r ResourceGetter, verb auth.Verb) func(c *gin.Context) {
	return func(c *gin.Context) {
		user, err := getUser(c)
		//subject, err := server.Auth.Subject(Token(c))
		if err != nil {
			c.HTML(http.StatusUnauthorized, "/errors/unauthorized.html", gin.H{})
			c.Abort()
			return
		}
		resource, err := r(c)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{})
			c.Abort()
			return
		}

		allowed, err := server.Auth.IsAllow(resource, verb, user)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{"err": err.Error()})
			c.Abort()
			return
		}
		if !allowed {
			logrus.Trace("rejected")
			c.HTML(http.StatusForbidden, "/errors/forbidden.html", gin.H{})
			c.Abort()
			return
		}
		logrus.Trace("accepted")

	}
}

func Token(c *gin.Context) *auth.Token {
	token, ok := c.Get(TOKEN)
	if ok {
		return token.(*auth.Token)
	}
	return nil
}
