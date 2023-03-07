package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (handler *Handler) Messages(c *gin.Context) {

	user, err := User(c)
	if err != nil {
		HTTPErrorHandler(c, err, WithAuthHeader)
		return
	}
	messages, err := handler.backend.GetMessages("user/" + user.Name)
	if err != nil {
		HTTPErrorHandler(c, err)
		return
	}
	logrus.Trace(messages)

	c.HTML(http.StatusOK, "/messages.html", gin.H{"messages": messages})
}
