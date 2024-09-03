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
				c.Error(HttpError{code: http.StatusNotFound, message: "md file not found"})
				return
			}
			renderedData, err := backend.Render(data)
			if err != nil {
				c.Error(err)
				return
			}

			footerData, err := backend.Plugin.GetWikiFooter(c.Request.URL.Path)
			if err != nil {
				logrus.Warn(err)
				c.Error(err)
				return
			}

			c.HTML(http.StatusOK, "notes/markdown.html", With(c, KeyValues{
				"content": template.HTML(renderedData),
				"footers": footerData,
			}))
			return
		}
	case "image":
		//handler.backend.FileExist(
		c.HTML(http.StatusOK, "notes/image.html", gin.H{
			"path": c.Request.URL.Path,
		})
	case "video":
		c.HTML(http.StatusOK, "notes/video.html", gin.H{
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
	logrus.Tracef("[View] serve raw file: '%s'", c.Request.URL.Path)
	backend := injector.Backends(c)
	// TODO serve partial content

	rs, info, err := backend.FileReadStream(c.Request.URL.Path)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	defer rs.Close()

	// github.com/gabriel-vasile/mimetype
	// mtype := mimetype.Detect(buf)
	http.ServeContent(c.Writer, c.Request, filepath.Base(c.Request.URL.Path), info.ModTime(), rs)
}
func EditForm(c *gin.Context) {
	backend := injector.Backends(c)

	category, _ := filetype(c.Request.URL.Path)

	switch category {
	case "text":
		data, err := backend.FileRead(c.Request.URL.Path)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
		c.HTML(http.StatusOK, "notes/editor.html", With(c, KeyValues{
			"data": template.HTML(data),
		}))
	default:
		c.HTML(http.StatusOK, "notes/upload.html", With(c, KeyValues{
			"path": c.Request.URL.Path,
		}))
	}
}
func UpdateWithForm(c *gin.Context) {
	backend := injector.Backends(c)

	p := c.Request.URL.Path
	req := &struct {
		Data string `form:"data"`
	}{}

	if err := c.ShouldBind(req); err != nil {
		c.Error(err)
		return
	}

	logrus.Tracef("data: %s", req.Data)

	if err := backend.FileWrite(p, []byte(req.Data)); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.Redirect(http.StatusSeeOther, p)
}
func MoveNote(c *gin.Context) {
	req := struct {
		Target string `form:"target"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if err := injector.Backends(c).FileMove(c.Request.URL.Path, req.Target); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Redirect(http.StatusSeeOther, req.Target)
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

	c.HTML(http.StatusOK, "notes/files.html", With(c, KeyValues{
		"files": files,
	}))
}
func DeleteNote(c *gin.Context) {
	backend := injector.Backends(c)

	if err := backend.FileDelete(c.Request.URL.Path); err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
}

// APIS

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

func Upload(c *gin.Context) {
	backend := injector.Backends(c)
	p := c.Request.URL.Path
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	if err := backend.FileWrite(p, data); err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.Status(http.StatusAccepted)
}
func Delete(c *gin.Context) {
	backend := injector.Backends(c)
	if c.GetHeader("X-Confirm") != path.Base(c.Request.URL.Path) {
		c.Error(HttpError{code: http.StatusBadRequest, message: "not confirmed"})
		c.Abort()
		return
	}

	err := backend.FileDelete(c.Request.URL.Path)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	c.Status(http.StatusNoContent)
}
