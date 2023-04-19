package handler

import (
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

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
func (handler *Handler) View(c *gin.Context) {
	logrus.Trace("view handler")
	category, subtype := filetype(c.Request.URL.Path)
	logrus.Tracef("category: %s, subtype: %s", category, subtype)

	switch category {
	case "text":
		switch subtype {
		case "markdown":
			data, err := handler.backend.FileRead(c.Request.URL.Path)
			if err != nil {
				logrus.Warnf("md file not found, %s", err)
				c.HTML(http.StatusNotFound, "/errors/not-found.html", gin.H{})
				return
			}
			renderedData, err := handler.backend.Render(data)
			if err != nil {
				c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{})
				return
			}

			footerData, err := handler.backend.Plugin.GetWikiFooter(c.Request.URL.Path)
			if err != nil {
				logrus.Warn(err)
				c.HTML(http.StatusInternalServerError, "/errors/internal-server-error.html", gin.H{})
				return
			}

			c.HTML(http.StatusOK, PageMarkdown, gin.H{
				"content": template.HTML(renderedData),
				"footers": footerData,
			})
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

func (handler *Handler) Raw(c *gin.Context) {
	logrus.Infof("[View] serve raw file: '%s'", c.Request.URL.Path)
	// TODO serve partial content

	rs, err := handler.backend.FileReadStream(c.Request.URL.Path)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// github.com/gabriel-vasile/mimetype
	// mtype := mimetype.Detect(buf)
	http.ServeContent(c.Writer, c.Request, "", time.Time{}, rs)
	//c.Data(http.StatusOK, http.DetectContentType(buf), buf)
}
func (handler *Handler) EditForm(c *gin.Context) {
	category, _ := filetype(c.Request.URL.Path)

	switch category {
	case "text":
		data, err := handler.backend.FileRead(c.Request.URL.Path)
		c.HTML(http.StatusOK, "/editor.html", gin.H{
			"data":  template.HTML(data),
			"path":  c.Request.URL.Path,
			"isNew": err != nil,
		})
	default:
		c.HTML(http.StatusOK, "/upload.html", gin.H{
			"path": c.Request.URL.Path,
		})
	}
}
func (handler *Handler) UpdateWithForm(c *gin.Context) {
	p := c.Request.URL.Path
	req := &struct {
		Data string `form:"data"`
	}{}

	if err := c.ShouldBind(req); err != nil {
		c.HTML(http.StatusBadRequest, PageErrBadRequest, gin.H{})
		return
	}

	logrus.Tracef("data: %s", req.Data)

	if err := handler.backend.FileWrite(p, []byte(req.Data)); err != nil {
		c.HTML(http.StatusInternalServerError, PageErrInternalServerError, gin.H{})
		c.Abort()
		return
	}
	c.Redirect(http.StatusSeeOther, p)
}
func (handler *Handler) Update(c *gin.Context) {
	p := c.Request.URL.Path
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err := handler.backend.FileWrite(p, data); err != nil {
		c.HTML(http.StatusInternalServerError, PageErrInternalServerError, gin.H{})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (handler *Handler) Files(c *gin.Context) {
	path := c.Request.URL.Path
	if strings.HasSuffix(path, ".md") {
		c.Redirect(http.StatusSeeOther, strings.TrimSuffix(path, ".md")+"?files")
		return
	}
	files, err := handler.backend.FileList(path)
	if err != nil && !os.IsNotExist(err) {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.HTML(http.StatusOK, "/viewers/files.html", gin.H{
		"files": files,
	})
}
func (handler *Handler) Delete(c *gin.Context) {
	if c.GetHeader("X-Confirm") != path.Base(c.Request.URL.Path) {
		c.HTML(http.StatusBadRequest, PageErrBadRequest, gin.H{})
		return
	}

	err := handler.backend.FileDelete(c.Request.URL.Path)
	if err != nil {
		c.HTML(http.StatusInternalServerError, PageErrInternalServerError, gin.H{})
		return
	}
	c.Status(http.StatusNoContent)
}
func (handler *Handler) Preview(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	renderedData, err := handler.backend.Render(data)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Data(http.StatusOK, "text/html", renderedData)
}
