package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) Search(c *gin.Context) {
	query := c.Query("q")

	result, err := handler.backend.FileSearch(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "/search.html", gin.H{"result": result})
}
