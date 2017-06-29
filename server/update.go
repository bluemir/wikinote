package server

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUpdateForm(c *gin.Context) {
}
func HandleUpdate(c *gin.Context) {
	p := c.Request.URL.Path
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	Backend(c).File().Write(p, data)
	c.JSON(http.StatusOK, gin.H{})
}
