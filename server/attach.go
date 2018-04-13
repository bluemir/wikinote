package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandleAttachForm(c *gin.Context) {
	p := c.Param("path")
	list, err := Backend(c).File().List(p)
	if err != nil {
		logrus.Warn(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "/attach.html", Data(c).
		Set("path", c.Param("path")).
		Set("files", list))

}
