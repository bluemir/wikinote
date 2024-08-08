package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	ContextKeyManager = "__auth_manager__"
	ContextKeyUser    = "__auth_user__"
)

func User(c *gin.Context) (*auth.User, error) {
	backend := injector.Backend(c)
	// 1. try to get user from context
	if u, ok := c.Get(ContextKeyUser); ok {
		logrus.Debug(u)
		if user, ok := u.(*auth.User); ok {
			return user, nil
		}
	}

	// 2. check basic auth or token

	user, err := backend.Auth.HTTP(c.Request)
	if err != nil {
		return nil, err
	}

	// set user to context
	c.Set(ContextKeyUser, user)

	return user, nil
}

type ResourceGetter func(c *gin.Context) (auth.Resource, error)

func Can(verb auth.Verb, getResource ResourceGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := User(c)
		if err != nil && !errors.Is(err, auth.ErrUnauthorized) {
			c.Error(err)
			c.Abort()
			return
		}
		// at this point, user can be nil(it means not logined user)

		resource, err := getResource(c)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		backend := injector.Backend(c)

		if err := backend.Auth.Can(user, verb, resource); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}
}

func (handler *Handler) CanAPI(c *gin.Context) {
	user, err := User(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	verb := c.Param("verb")
	kind := strings.TrimPrefix(c.Param("kind"), "/")

	logrus.WithField("verb", verb).WithField("kind", kind).Trace("API called")

	if err := handler.backend.Auth.Can(user, auth.Verb(verb), auth.KeyValues{"kind": kind}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusOK)
}
