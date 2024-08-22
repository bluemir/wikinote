package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Messages(c *gin.Context) {
	backend := injector.Backends(c)

	user, err := User(c)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	messages, err := backend.GetMessages("user/" + user.Name)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	logrus.Trace(messages)

	c.HTML(http.StatusOK, "/messages.html", gin.H{"messages": messages})
}
