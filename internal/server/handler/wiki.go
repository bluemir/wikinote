package handler

import (
	"html/template"
	"io"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/bluemir/wikinote/internal/server/injector"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func parseMIME(mime string) (string, string) {
	arr := strings.SplitN(strings.SplitN(mime, ";", 2)[0], "/", 2)
	if len(arr) > 1 {
		return arr[0], arr[1]
	}
	return arr[0], ""
}
func filetype(path string) (string, string) {
	ctype := mime.TypeByExtension(filepath.Ext(path))

	if ctype != "" {
		return parseMIME(ctype)
	}
	// failback
	switch filepath.Ext(path) {
	case ".md":
		return "text", "markdown"
	default:
		return "", ""
	}
}
func View(c *gin.Context) {

	backend := injector.Backends(c)

	logrus.Trace("view handler")
	category, subtype := filetype(c.Request.URL.Path)
	logrus.Tracef("category: %s, subtype: %s", category, subtype)

	switch category {
	case "text":
		switch subtype {
		case "markdown":
			data, err := backend.FileRead(c.Request.URL.Path)
			if err != nil {
				logrus.Warnf("md file not found, %s", err)
				c.HTML(http.StatusNotFound, PageErrNotFound, gin.H{})
				return
			}
			renderedData, err := backend.Render(data)
			if err != nil {
				c.Error(err)
				c.HTML(http.StatusInternalServerError, PageErrInternalServerError, gin.H{})
				return
			}

			footerData, err := backend.Plugin.GetWikiFooter(c.Request.URL.Path)
			if err != nil {
				logrus.Warn(err)
				c.Error(err)
				c.HTML(http.StatusInternalServerError, PageErrInternalServerError, gin.H{})
				return
			}

			//c.HTML(http.StatusOK, PageMarkdown, gin.H{
			//	"content": template.HTML(renderedData),
			//	"footers": footerData,
			//})
			c.HTML(http.StatusOK, PageMarkdown, With(c, KeyValues{
				"content": template.HTML(renderedData),
				"footers": footerData,
			}))
			return
		}
	case "image":
		//handler.backend.FileExist(
		c.HTML(http.StatusOK, PageImage, gin.H{
			"path": c.Request.URL.Path,
		})
	case "video":
		c.HTML(http.StatusOK, PageVideo, gin.H{
			"path": c.Request.URL.Path,
		})
	default:
		if !strings.HasSuffix(c.Request.URL.Path, ".md") {
			// or show files?
			c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path+".md")
		}
	}
}

func Raw(c *gin.Context) {
	logrus.Infof("[View] serve raw file: '%s'", c.Request.URL.Path)
	backend := injector.Backends(c)
	// TODO serve partial content

	rs, err := backend.FileReadStream(c.Request.URL.Path)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	defer rs.Close()

	// github.com/gabriel-vasile/mimetype
	// mtype := mimetype.Detect(buf)
	http.ServeContent(c.Writer, c.Request, "", time.Time{}, rs)
	//c.Data(http.StatusOK, http.DetectContentType(buf), buf)
}
func EditForm(c *gin.Context) {
	backend := injector.Backends(c)

	category, _ := filetype(c.Request.URL.Path)

	switch category {
	case "text":
		data, err := backend.FileRead(c.Request.URL.Path)
		c.HTML(http.StatusOK, PageEditor, gin.H{
			"data":  template.HTML(data),
			"path":  c.Request.URL.Path,
			"isNew": err != nil,
		})
	default:
		c.HTML(http.StatusOK, PageUpload, gin.H{
			"path": c.Request.URL.Path,
		})
	}
}
func UpdateWithForm(c *gin.Context) {
	backend := injector.Backends(c)

	p := c.Request.URL.Path
	req := &struct {
		Data string `form:"data"`
	}{}

	if err := c.ShouldBind(req); err != nil {
		c.HTML(http.StatusBadRequest, PageErrBadRequest, gin.H{})
		return
	}

	logrus.Tracef("data: %s", req.Data)

	if err := backend.FileWrite(p, []byte(req.Data)); err != nil {
		c.HTML(http.StatusInternalServerError, PageErrInternalServerError, gin.H{})
		c.Abort()
		return
	}
	c.Redirect(http.StatusSeeOther, p)
}
func Update(c *gin.Context) {
	backend := injector.Backends(c)
	p := c.Request.URL.Path
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := backend.FileWrite(p, data); err != nil {
		c.HTML(http.StatusInternalServerError, PageErrInternalServerError, gin.H{})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func Files(c *gin.Context) {
	backend := injector.Backends(c)
	path := c.Request.URL.Path
	if strings.HasSuffix(path, ".md") {
		c.Redirect(http.StatusSeeOther, strings.TrimSuffix(path, ".md")+"?files")
		return
	}
	files, err := backend.FileList(path)
	if err != nil && !os.IsNotExist(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, PageFiles, With(c, KeyValues{
		"files": files,
	}))
}
func Delete(c *gin.Context) {
	backend := injector.Backends(c)
	if c.GetHeader("X-Confirm") != path.Base(c.Request.URL.Path) {
		c.HTML(http.StatusBadRequest, PageErrBadRequest, gin.H{})
		return
	}

	err := backend.FileDelete(c.Request.URL.Path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, PageErrInternalServerError, gin.H{})
		return
	}
	c.Status(http.StatusNoContent)
}
func Preview(c *gin.Context) {
	backend := injector.Backends(c)
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	renderedData, err := backend.Render(data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Data(http.StatusOK, "text/html", renderedData)
}
