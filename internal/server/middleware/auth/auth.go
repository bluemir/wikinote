package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/auth"
)

type Resource = auth.Resource
type Verb = auth.Verb
type KeyValues = auth.KeyValues

const (
	ContextKeyManager = "__auth_manager__"
	ContextKeyUser    = "__auth_user__"
)

func Middleware(m *auth.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Attach manager
		c.Set(ContextKeyManager, m)

		// try to login
		user, err := m.HTTP(c.Request)
		if err != nil {
			return
		}
		if user == nil {
			return
		}
		c.Set(ContextKeyUser, user)
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

type ResourceGetter func(c *gin.Context) (auth.Resource, error)

func Authz(getResource ResourceGetter, verb Verb) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := User(c)
		if err != nil && !errors.Is(err, auth.ErrUnauthorized) {
			c.Error(err)
			c.Abort()
			return
		}

		resource, err := getResource(c)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if err := manager(c).Can(user, verb, resource); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}
}
