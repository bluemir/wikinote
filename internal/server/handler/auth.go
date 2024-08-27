package handler

import (
	"encoding/gob"
	"errors"
	"net/http"
	"strings"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type typeSesstionKeyUser struct{}

var SessionKeyUser = typeSesstionKeyUser{}

const (
	ContextKeyManager = "__auth_manager__"
	ContextKeyUser    = "__auth_user__"
)

func init() {
	gob.Register(&auth.User{})
	gob.Register(typeSesstionKeyUser{})
}

func User(c *gin.Context) (*auth.User, error) {
	backend := injector.Backends(c)
	// 1. try to get user from context
	if u, ok := c.Get(ContextKeyUser); ok {
		if user, ok := u.(*auth.User); ok {
			return user, nil
		}
		logrus.Tracef("found user in context but not matched type: %T", u)
	}

	logrus.Trace("user not found in context")

	// 2. check session
	session := sessions.Default(c)
	u := session.Get(SessionKeyUser)
	if u != nil {
		if user, ok := u.(*auth.User); ok {
			c.Set(ContextKeyUser, user)
			return user, nil
		}
		logrus.Tracef("found user in session but not matched type: %T", u)
	}

	logrus.Trace("user not found in session")

	// 3. check basic auth or token
	user, err := backend.Auth.HTTP(c.Request)
	if err != nil {
		return nil, err
	}
	if u != nil {
		// set user to context
		c.Set(ContextKeyUser, user)
	}

	return user, nil
}

func Login(c *gin.Context) {
	req := struct {
		Username string `form:"username"`
		Password string `form:"password"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	user, err := injector.Backends(c).Auth.Default(req.Username, req.Password)
	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}

	session := sessions.Default(c)
	session.Set(SessionKeyUser, user)

	if err := session.Save(); err != nil {
		c.Error(err)
		return
	}

	logrus.Tracef("logined: %s", user.Name)

	c.Redirect(http.StatusSeeOther, returnURL(c))
}
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(SessionKeyUser)
	session.Save()

	c.Redirect(http.StatusSeeOther, returnURL(c))
}

func returnURL(c *gin.Context) string {
	returnURL := c.Query("return")
	if returnURL == "" {
		returnURL = "/"
	}

	if !strings.HasPrefix(returnURL, "/") {
		returnURL = "/" // Must start '/', not 'https://', 'javascript:', ...
	}

	return returnURL
}

type ResourceGetter func(c *gin.Context) (auth.Resource, error)

func Can(verb auth.Verb, getResource ResourceGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		resource, err := getResource(c)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		user, err := User(c)
		if err != nil && !errors.Is(err, auth.ErrUnauthorized) {
			c.Error(err)
			c.Abort()
			return
		}
		// at this point, user can be nil(it means not logined user)

		backend := injector.Backends(c)

		if err := backend.Auth.Can(user, verb, resource); err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}
}

func CanAPI(c *gin.Context) {
	backend := injector.Backends(c)
	user, err := User(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	verb := c.Param("verb")
	kind := strings.TrimPrefix(c.Param("kind"), "/")

	logrus.WithField("verb", verb).WithField("kind", kind).Trace("API called")

	if err := backend.Auth.Can(user, auth.Verb(verb), auth.KeyValues{"kind": kind}); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusOK)
}
