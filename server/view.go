package server

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandleView(c *gin.Context) {
	//logrus.Debugf("path: %s, accept: %+v", c.Request.URL.Path, c.GetHeader("Accept"))

	// if raw exist
	if _, ok := c.GetQuery("raw"); ok {
		logrus.Infof("[View] serve raw file: '%s'", Backend(c).File().GetFullPath(c.Request.URL.Path))
		c.File(Backend(c).File().GetFullPath(c.Request.URL.Path))
		return
	}

	// plugin renderer?
	switch {
	// TODO check file type
	case checkExt(c, ".md"):
		data, err := Backend(c).File().Read(c.Request.URL.Path)
		if err != nil {
			logrus.Warnf("md file not found, %s", err)
			c.HTML(http.StatusNotFound, "/errors/not-found.html", Data(c))
			return
		}
		renderedData, err := Backend(c).Render(data)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "/view/internal-error.html", Data(c))
			return
		}

		// error will be
		c.HTML(http.StatusOK, "/view/markdown.html", Data(c).
			Set("data", template.HTML(renderedData)).
			Set("footer", Backend(c).Plugin().Footer(c.Request.URL.Path)),
		)
		// markdown
	case checkExt(c, ".jpg", ".png", ".bmp", ".gif"):
		// image
		// check exist?
		c.HTML(http.StatusOK, "/view/image.html", Data(c))
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
