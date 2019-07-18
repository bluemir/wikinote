package server

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkgs/auth"
	"github.com/bluemir/wikinote/pkgs/renderer"
)

func HandleView(c *gin.Context) {
	//logrus.Debugf("path: %s, accept: %+v", c.Request.URL.Path, c.GetHeader("Accept"))

	switch {
	// TODO check file type
	case checkExt(c, ".md"):
		data, err := Backend(c).File().Read(c.Request.URL.Path)
		if err != nil {
			logrus.Warnf("md file not found, %s", err)
			c.HTML(http.StatusNotFound, "/errors/not-found.html", renderer.Data{}.With(c))
			return
		}
		authCtx, ok := c.Get(AUTH_CONTEXT)
		if !ok {
			c.HTML(http.StatusInternalServerError, "/view/internal-error.html", renderer.Data{}.With(c))
			c.Abort()
			return
		}
		buf, err := Backend(c).Plugin().OnReadWiki(authCtx.(*auth.Context), c.Request.URL.Path, data)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "/view/internal-error.html", renderer.Data{}.With(c))
			return
		}

		renderedData, err := Backend(c).Render(buf)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "/view/internal-error.html", renderer.Data{}.With(c))
			return
		}

		//renderer.Of(c).AddFlashInfo(sdaasd)
		/*c.HTML(http.StatusOK, "/view/markdown.html", &renderer.Data{
			"data": template.HTML(renderedData),
			"footer": Backend(c).Plugin().Fotter(c.Request.URL.Path)),
		}.With(c))*/
		c.HTML(http.StatusOK, "/view/markdown.html", renderer.Data{
			"data":   template.HTML(renderedData),
			"footer": Backend(c).Plugin().Footer(c.Request.URL.Path),
		}.With(c))
		// markdown
	case checkExt(c, ".jpg", ".png", ".bmp", ".gif"):
		// image
		// check exist?
		c.HTML(http.StatusOK, "/view/image.html", renderer.Data{}.With(c))
	case checkExt(c, ".mp4"):
		// video
	case checkExt(c, ".mp3"):
		// music
	default:
		// check ext and
		// render md files or other rendering
		logrus.Debugf("no ext: %s", c.Request.URL.Path)
		// TODO CONFIG
		c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path+".md")
	}
}
func checkExt(c *gin.Context, ext ...string) bool {
	for _, e := range ext {
		if ori := strings.ToLower(filepath.Ext(c.Request.URL.Path)); ori == e {
			return true
		}
	}
	return false
}
func HandleRaw(c *gin.Context) {
	logrus.Infof("[View] serve raw file: '%s'", Backend(c).File().GetFullPath(c.Request.URL.Path))
	c.File(Backend(c).File().GetFullPath(c.Request.URL.Path))
	return
}
func HandleEditForm(c *gin.Context) {
	backend := Backend(c)
	path := c.Request.URL.Path

	data, err := backend.File().Read(path)
	if err != nil {
		//c.AbortWithError(http.StatusNotFound, err)
		renderer.Of(c).Info("Create new note")
		// TODO add flash message
	}

	c.HTML(http.StatusOK, "/edit.html", renderer.Data{
		"data": string(data),
		"path": c.Param("path"),
	}.With(c))
}
func HandleUpdateForm(c *gin.Context) {
	p := c.Request.URL.Path
	data, ok := c.GetPostForm("data")
	if !ok {
		c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error on form"))
	}
	err := Backend(c).File().Write(p, []byte(data))
	if err != nil {

		c.HTML(http.StatusInternalServerError, "/errors/not-found.html", gin.H{})
		c.Abort()
	}
	c.Redirect(http.StatusSeeOther, p)
}
func HandleUpdate(c *gin.Context) {
	p := c.Request.URL.Path
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	Backend(c).File().Write(p, data)
	c.JSON(http.StatusOK, gin.H{})
}
func HandleAttachForm(c *gin.Context) {
	p := c.Param("path")
	list, err := Backend(c).File().List(p)
	if err != nil {
		logrus.Warn(err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.HTML(http.StatusOK, "/attach.html", renderer.Data{
		"path":  c.Param("path"),
		"files": list,
	}.With(c))
}
