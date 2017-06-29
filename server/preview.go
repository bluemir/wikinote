package server

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlePreview(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	renderedData, err := Backend(c).Render(data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Data(http.StatusOK, "text/html", renderedData)
}
