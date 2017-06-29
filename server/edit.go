package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandleEditForm(c *gin.Context) {
	backend := Backend(c)
	logrus.Infof("info:%s", c.Param("path"))
	data, err := backend.File().Read(c.Param("path"))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.HTML(http.StatusOK, "/edit.html", Data(c).
		Set("data", string(data)).
		Set("path", c.Param("path")),
	)
}
