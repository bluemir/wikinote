package server

import (
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func parseMIME(mime string) (string, string) {
	arr := strings.SplitN(strings.SplitN(mime, ";", 2)[0], "/", 2)
	return arr[0], arr[1]
}
func (server *Server) HandleView(c *gin.Context) {
	//logrus.Debugf("path: %s, accept: %+v", c.Request.URL.Path, c.GetHeader("Accept"))

	//mime.TypeByExtension
	// mime
	ctype := mime.TypeByExtension(filepath.Ext(c.Request.URL.Path))
	category, subtype := parseMIME(ctype)
	logrus.Debugf("mime: %s", ctype)

	switch category {
	case "text":
		switch subtype {
		case "markdown":
			data, err := server.FileRead(c.Request.URL.Path)
			if err != nil {
				logrus.Warnf("md file not found, %s", err)
				c.HTML(http.StatusNotFound, "/errors/not-found.html", gin.H{})
				return
			}
			renderedData, err := server.Render(data)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{})
				return
			}

			footerData, err := server.Backend.Plugin.WikiFooter(c.Request.URL.Path)
			if err != nil {
				logrus.Warn(err)
				c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{})
				return
			}

			c.HTML(http.StatusOK, "/views/markdown.html", gin.H{
				"data":       template.HTML(renderedData),
				"footerData": footerData,
			})
			return
		}
		//case "video":
	case "image":
		c.HTML(http.StatusOK, "/view/image.html", c)
	default:
		c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path+".md")
	}
}

func (server *Server) HandleRaw(c *gin.Context) {
	logrus.Infof("[View] serve raw file: '%s'", c.Request.URL.Path)

	buf, err := server.FileRead(c.Request.URL.Path)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// github.com/gabriel-vasile/mimetype
	// mtype := mimetype.Detect(buf)

	c.Data(http.StatusOK, http.DetectContentType(buf), buf)
}
func (server *Server) HandleEditForm(c *gin.Context) {
	path := c.Request.URL.Path

	data, err := server.FileRead(path)
	if err != nil {
		// Create new file
		// c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{"msg": err.Error()})
		// c.Abort()
		// return
	}

	c.HTML(http.StatusOK, "/editor.html", gin.H{
		"data": string(data),
		"path": c.Request.URL.Path,
	})
}
func (server *Server) HandleUpdateWithForm(c *gin.Context) {
	p := c.Request.URL.Path
	req := &struct {
		Data string `form:"data"`
	}{}

	if err := c.ShouldBind(req); err != nil {
		c.HTML(http.StatusBadRequest, "/errors/internal-error.html", gin.H{})
		return
	}

	if err := server.FileWrite(p, []byte(req.Data)); err != nil {
		c.HTML(http.StatusInternalServerError, "/errors/not-found.html", gin.H{})
		c.Abort()
		return
	}
	c.Redirect(http.StatusSeeOther, p)
}
func (server *Server) HandleUpdate(c *gin.Context) {
	p := c.Request.URL.Path
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := server.FileWrite(p, data); err != nil {
		c.HTML(http.StatusInternalServerError, "/errors/not-found.html", gin.H{})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
func (server *Server) HandleDeleteForm(c *gin.Context) {
	c.HTML(http.StatusOK, "/delete.html", gin.H{
		"name": path.Base(c.Request.URL.Path),
	})
}
func (server *Server) HandleDelete(c *gin.Context) {
	if c.GetHeader("X-Confirm") != path.Base(c.Request.URL.Path) {
		c.HTML(http.StatusBadRequest, "/errors/bad-request.html", gin.H{})
		return
	}

	err := server.Backend.FileDelete(c.Request.URL.Path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "/errors/internal-sever-error.html", gin.H{})
		return
	}
	c.Status(http.StatusNoContent)
}
func (server *Server) HandlePreview(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	renderedData, err := server.Render(data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Data(http.StatusOK, "text/html", renderedData)
}
