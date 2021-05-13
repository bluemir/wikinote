package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) HandleSearch(c *gin.Context) {
	query := c.Query("q")

	result, err := server.FileSearch(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "/search.html", gin.H{"result": result})
}
