package server

import (
	"github.com/bluemir/wikinote/backend/user"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionContext interface {
	Login(user *user.User)
	Logout()

	Username() string
	Role() string
}

func Session(c *gin.Context) SessionContext {
	session := sessions.Default(c)

	return &sessionContext{
		session,
	}
}

type sessionContext struct {
	sessions.Session
}

func (sc *sessionContext) Login(user *user.User) {
	sc.Set(USERNAME, user.Id)
	sc.Set(ROLE, user.Role)
	sc.Save()
}
func (sc *sessionContext) Logout() {
	sc.Delete(USERNAME)
	sc.Delete(ROLE)
	sc.Save()
}
func (sc *sessionContext) Username() string {
	if username, ok := sc.Get(USERNAME).(string); ok {
		return username
	} else {
		return ""
	}
}
func (sc *sessionContext) Role() string {
	if role, ok := sc.Get(ROLE).(string); ok {
		return role
	} else {
		return ""
	}
}
