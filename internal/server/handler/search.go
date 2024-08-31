package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) {
	backend := injector.Backends(c)
	query := c.Query("q")

	result, err := backend.FileSearch(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, PageSearch, With(c, result))
}
