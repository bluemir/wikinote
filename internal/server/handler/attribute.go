package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/fileattr"
)

func (handler *Handler) AttributeGet(c *gin.Context) {
	attrs, err := handler.backend.FileAttr.Find(&fileattr.FileAttr{Path: c.Request.URL.Path})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	result := map[string]string{}
	for _, kv := range attrs {
		result[kv.Key] = kv.Value
	}
	c.JSON(http.StatusOK, result)
}
func (handler *Handler) AttributeUpdate(c *gin.Context) {
	req := map[string]string{}
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "/errors/internal-server-error.html", gin.H{})
		return
	}
	for k, v := range req {
		if err := handler.backend.FileAttr.Save(&fileattr.FileAttr{
			Path:  c.Request.URL.Path,
			Key:   k,
			Value: v,
		}); err != nil {
			c.HTML(http.StatusInternalServerError, "/errors/not-found.html", gin.H{})
			c.Abort()
			return
		}
	}
	attrs, err := handler.backend.FileAttr.Find(&fileattr.FileAttr{
		Path: c.Request.URL.Path,
	})
	if err != nil {
		c.HTML(http.StatusInternalServerError, "/errors/not-found.html", gin.H{})
		c.Abort()
		return
	}
	for _, attr := range attrs {
		if _, ok := req[attr.Key]; ok {
			continue
		}

		if err := handler.backend.FileAttr.Delete(&attr); err != nil {
			c.HTML(http.StatusInternalServerError, "/errors/not-found.html", gin.H{})
			c.Abort()
			return
		}
	}

	c.JSON(http.StatusOK, req)
}
