package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/pkgs/renderer"
)

func (server *Server) HandleSearch(c *gin.Context) {
	query := c.Query("q")

	// TODO html
	result, err := server.File().Search(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "/search.html", renderer.Data{
		"result": result,
	}.With(c))
}
