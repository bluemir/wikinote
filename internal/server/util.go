package server

import (
	"crypto"
	"encoding/hex"
	"io"
	"net/http"
	"strings"

	"github.com/bluemir/wikinote/internal/buildinfo"
	"github.com/gin-gonic/gin"
)

func (server *Server) static(path string) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, path, c)
	}
}

func (server *Server) initEtag() error {
	hashed := crypto.SHA512.New()

	io.WriteString(hashed, buildinfo.AppName)
	io.WriteString(hashed, buildinfo.Version)
	io.WriteString(hashed, buildinfo.BuildTime)

	server.etag = hex.EncodeToString(hashed.Sum(nil))[:20]

	return nil
}

func (server *Server) staticCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, max-age=86400")
	c.Header("ETag", server.etag)

	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, server.etag) {
			c.Status(http.StatusNotModified)
			c.Abort()
			return
		}
	}

	c.Request.Header.Del("If-Modified-Since") // only accept etag
}
