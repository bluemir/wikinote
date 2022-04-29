package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/handler"
)

type Resource = auth.Resource
type Verb = auth.Verb
type KeyValues = auth.KeyValues

const (
	ContextKeyManager = "__auth_manager__"
	ContextKeyUser    = "__auth_user__"

	realm = "wikinote.bluemir.me"
)

func Middleware(m *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(ContextKeyManager, m)
	}
}
func manager(c *gin.Context) *auth.Manager {
	return c.MustGet(ContextKeyManager).(*auth.Manager)
}

func User(c *gin.Context) (*auth.User, error) {
	// 1. try to get user from context
	if u, ok := c.Get(ContextKeyUser); ok {
		logrus.Debug(u)
		if user, ok := u.(*auth.User); ok {
			return user, nil
		}
	}

	// 2. check basic auth or token
	user, err := manager(c).HTTP(c.Request)
	if err != nil {
		return nil, err
	}

	// set user to context
	c.Set(ContextKeyUser, user)

	return user, nil
}

func RequestLogin(c *gin.Context) {
	_, err := User(c)

	if errors.Is(err, auth.ErrUnauthorized) {
		handler.ErrorHandler(c, err)
		return
	}
}

type ResourceGetter func(c *gin.Context) (auth.Resource, error)

func Authz(getResource ResourceGetter, verb Verb) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := User(c)
		if err != nil && !errors.Is(err, auth.ErrUnauthorized) {
			handler.ErrorHandler(c, err)
			return
		}

		resource, err := getResource(c)
		if err != nil {
			handler.ErrorHandler(c, err)
			return
		}

		if err := manager(c).IsAllow(resource, verb, user); err != nil {
			handler.ErrorHandler(c, err)
			return
		}
	}
}
func IssueToken(c *gin.Context) {
	req := struct {
		Username string
		Password string
	}{}

	if err := c.ShouldBind(&req); err != nil {
		handler.ErrorHandler(c, err)
		return
	}

	user, err := manager(c).Default(req.Username, req.Password)
	if err != nil {
		handler.ErrorHandler(c, err)
		return
	}

	t := time.Now().Add(6 * time.Hour)
	token, err := manager(c).NewHTTPToken(user.Name, t)
	if err != nil {
		handler.ErrorHandler(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"expiredAt": t.Format(time.RFC3339),
	})
}
func RevokeToken(c *gin.Context) {
	if err := manager(c).RevokeHTTPToken(c.Request); err != nil {
		handler.ErrorHandler(c, err)
		return
	}
}
