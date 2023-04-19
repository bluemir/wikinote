package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) Admin(c *gin.Context) {
	c.HTML(http.StatusOK, "/admin/index.html", c)
}
