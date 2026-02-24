package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ListAllMessages(c *gin.Context) error {
	backend := injector.Backends(c)

	messages, err := backend.GetEvents(c.Request.Context())
	if err != nil {
		return err
	}
	logrus.Trace(messages)

	c.HTML(http.StatusOK, "admin/messages.html", With(c, gin.H{"messages": messages}))
	return nil
}
