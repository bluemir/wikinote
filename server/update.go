package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleUpdateForm(c *gin.Context) {
	p := c.Request.URL.Path
	data, ok := c.GetPostForm("data")
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error on form"))
	}
	err := Backend(c).File().Write(p, []byte(data))
	if err != nil {

		c.HTML(http.StatusInternalServerError, "/errors/not-found.html", gin.H{})
		c.Abort()
	}
	c.Redirect(http.StatusSeeOther, p)
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
