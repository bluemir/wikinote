package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleAttachForm(c *gin.Context) {
	c.HTML(http.StatusOK, "/attach.html", Data(c).
		Set("path", c.Param("path")),
	)
}
