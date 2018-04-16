package server

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type SessionContext interface {
}

func Session(c *gin.Context) SessionContext {
	return sessions.Default(c)
}
