package handler

import (
	"net/http"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
)

func Search(c *gin.Context) error {
	backend := injector.Backends(c)
	query := c.Query("q")

	result, err := backend.FileSearch(query)
	if err != nil {
		return err
	}
	c.HTML(http.StatusOK, "notes/search.html", With(c, result))
	return nil
}
