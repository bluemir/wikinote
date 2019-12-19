package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/pkg/fileattr"
)

func (server *Server) HandleAttributeGet(c *gin.Context) {
	attrs, err := server.Backend.FileAttr.Find(&fileattr.FileAttr{Path: c.Request.URL.Path})
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	result := map[string]string{}
	for _, kv := range attrs {
		result[kv.Key] = kv.Value
	}
	c.JSON(http.StatusOK, result)
}
func (server *Server) HandleAttributeUpdate(c *gin.Context) {
	req := map[string]string{}
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "/errors/internal-server-error.html", gin.H{})
		return
	}
	for k, v := range req {
		if err := server.Backend.FileAttr.Save(&fileattr.FileAttr{
			Path:  c.Request.URL.Path,
			Key:   k,
			Value: v,
		}); err != nil {
			c.HTML(http.StatusInternalServerError, "/errors/not-found.html", gin.H{})
			c.Abort()
			return
		}
	}
	attrs, err := server.Backend.FileAttr.Find(&fileattr.FileAttr{
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

		if err := server.Backend.FileAttr.Delete(&attr); err != nil {
			c.HTML(http.StatusInternalServerError, "/errors/not-found.html", gin.H{})
			c.Abort()
			return
		}
	}

	c.JSON(http.StatusOK, req)
}
