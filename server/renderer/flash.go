package renderer

import (
	"fmt"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	MSG_WARN = "WARN"
	MSG_INFO = "INFO"
)

type MessageContext interface {
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
}

func NewMessageContext(c *gin.Context) MessageContext {
	return &msgContext{
		session: sessions.Default(c),
	}
}

type msgContext struct {
	session sessions.Session
}

func (msg *msgContext) Info(format string, v ...interface{}) {
	logrus.Infof(format, v...)
	msg.session.AddFlash(fmt.Sprintf(format, v...), MSG_INFO)
	msg.session.Save()
}
func (msg *msgContext) Warn(format string, v ...interface{}) {
	logrus.Warnf(format, v...)
	msg.session.AddFlash(fmt.Sprintf(format, v...), MSG_WARN)
	msg.session.Save()
}
