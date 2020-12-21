package server

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (server *Server) HandleView(c *gin.Context) {
	//logrus.Debugf("path: %s, accept: %+v", c.Request.URL.Path, c.GetHeader("Accept"))

	switch strings.ToLower(filepath.Ext(c.Request.URL.Path)) {
	case ".md":
		data, err := server.FileRead(c.Request.URL.Path)
		if err != nil {
			logrus.Warnf("md file not found, %s", err)
			c.HTML(http.StatusNotFound, "/errors/not-found.html", gin.H{
				"breadcrumb": makeBreadcurmb(c.Request.URL.Path),
			})
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
			"breadcrumb": makeBreadcurmb(c.Request.URL.Path),
			"footerData": footerData,
		})
		// markdown
	case ".jpg", ".png", ".bmp", ".gif":
		// image
		// check exist?
		c.HTML(http.StatusOK, "/view/image.html", gin.H{})
	case ".mp4":
		// video
	case ".mp3":
		// music
	default:
		// check ext and
		// render md files or other rendering
		logrus.Debugf("no ext: %s", c.Request.URL.Path)
		// TODO CONFIG
		c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path+".md")
	}
}

func (server *Server) HandleRaw(c *gin.Context) {
	logrus.Infof("[View] serve raw file: '%s'", server.GetFullPath(c.Request.URL.Path))
	c.File(server.GetFullPath(c.Request.URL.Path))
	return
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

	c.HTML(http.StatusOK, "/edit.html", gin.H{
		"data":       string(data),
		"path":       c.Request.URL.Path,
		"breadcrumb": makeBreadcurmb(c.Request.URL.Path),
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

// using in template
type Breadcrumb struct {
	Name string
	Path string
}

func makeBreadcurmb(p string) []Breadcrumb {
	result := []Breadcrumb{}
	arr := strings.Split(p, "/")
	for index, name := range arr {
		if name == "" {
			continue
		}
		p := path.Join(arr[:index+1]...)
		result = append(result, Breadcrumb{
			Name: name,
			Path: path.Join("/", p),
		})
	}
	return result
}
