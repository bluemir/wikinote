package handler

import (
	"html/template"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
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
	case "":
		return "directory", ""
	default:
		return "unknown", ""
	}
}
func View(c *gin.Context) error {

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
				return HttpError{code: http.StatusNotFound, message: "md file not found"}
			}
			renderedData, err := backend.Render(data)
			if err != nil {
				return err
			}

			footerData, err := backend.Plugin.GetFooters(c.Request.URL.Path)
			if err != nil {
				logrus.Warn(err)
				return err
			}

			c.HTML(http.StatusOK, "notes/markdown.html", With(c, KeyValues{
				"content": template.HTML(renderedData),
				"footers": footerData,
			}))
			return nil
		}
	case "image":
		//handler.backend.FileExist(
		c.HTML(http.StatusOK, "notes/image.html", With(c, gin.H{
			"path": c.Request.URL.Path,
		}))
	case "video":
		c.HTML(http.StatusOK, "notes/video.html", With(c, gin.H{
			"path": c.Request.URL.Path,
		}))
	case "directory":
		if !strings.HasSuffix(c.Request.URL.Path, ".md") {
			// or show files?
			c.Redirect(http.StatusTemporaryRedirect, c.Request.URL.Path+".md")
		}
	default:
		c.HTML(http.StatusOK, "notes/not-supported.html", With(c, nil))
	}
	return nil
}

func Raw(c *gin.Context) error {
	logrus.Tracef("[View] serve raw file: '%s'", c.Request.URL.Path)
	backend := injector.Backends(c)
	// TODO serve partial content

	rs, info, err := backend.FileReadStream(c.Request.URL.Path)
	if err != nil {
		return err
	}
	defer rs.Close()

	// github.com/gabriel-vasile/mimetype
	// mtype := mimetype.Detect(buf)
	http.ServeContent(c.Writer, c.Request, filepath.Base(c.Request.URL.Path), info.ModTime(), rs)
	return nil
}
func EditForm(c *gin.Context) error {
	backend := injector.Backends(c)

	kind, subkind := filetype(c.Request.URL.Path)
	logrus.Tracef("type=%s, subtype=%s", kind, subkind)

	switch kind {
	case "text":
		data, err := backend.FileRead(c.Request.URL.Path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		c.HTML(http.StatusOK, "editors/text.html", With(c, KeyValues{
			"data": string(data),
		}))
		return nil
	case "directory":
		// it may be directory. show upload form.
		files, err := backend.FileList(c.Request.URL.Path)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		c.HTML(http.StatusOK, "editors/files.html", With(c, KeyValues{
			"files": files,
		}))
		return nil
	case "unknown":
		fallthrough
	default:
		// it's file but there is no editor for this file
		// show upload forms for replacement
		c.HTML(http.StatusOK, "editors/upload.html", With(c, KeyValues{
			"path": c.Request.URL.Path,
		}))
		return nil
	}
}
func UpdateWithForm(c *gin.Context) error {
	backend := injector.Backends(c)
	// not work with file upload

	p := c.Request.URL.Path
	req := &struct {
		Data string `form:"data"`
		File *multipart.FileHeader
	}{}

	if err := c.ShouldBind(req); err != nil {
		return err
	}

	logrus.Tracef("data: %s", req.Data)

	if err := backend.FileWrite(p, []byte(req.Data)); err != nil {
		return err
	}
	c.Redirect(http.StatusSeeOther, p)
	return nil
}

func UploadFileToReplace(c *gin.Context) error {
	// replace file
	category, _ := filetype(c.Request.URL.Path)
	if category == "directory" {
		return ErrNotFile
	}
	backend := injector.Backends(c)

	req := &struct {
		File *multipart.FileHeader `form:"file"`
	}{}

	if err := c.ShouldBind(req); err != nil {
		return err
	}

	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if err := backend.FileWriteStream(c.Request.URL.Path, file); err != nil {
		return err
	}
	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
	return nil
	//c.SaveUploadedFile
}
func UploadFiles(c *gin.Context) error {
	// upload file to directory
	category, _ := filetype(c.Request.URL.Path)
	if category != "directory" {
		return ErrNotDirectory
	}

	backend := injector.Backends(c)

	req := &struct {
		File     *multipart.FileHeader `form:"file"`
		FileName string                `form:"filename"`
	}{}

	if err := c.ShouldBind(req); err != nil {
		return err
	}

	path := filepath.Join(c.Request.URL.Path, req.FileName)

	file, err := req.File.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if err := backend.FileWriteStream(path, file); err != nil {
		return err
	}
	c.Redirect(http.StatusSeeOther, c.Request.URL.Path+"?edit")
	return nil
}

func MoveNote(c *gin.Context) error {
	req := struct {
		Target string `form:"target"`
	}{}

	if err := c.ShouldBind(&req); err != nil {
		return err
	}

	if err := injector.Backends(c).FileMove(c.Request.URL.Path, req.Target); err != nil {
		return err
	}

	c.Redirect(http.StatusSeeOther, req.Target)
	return nil
}

func DeleteNote(c *gin.Context) error {
	backend := injector.Backends(c)

	if err := backend.FileDelete(c.Request.URL.Path); err != nil {
		return err
	}
	c.Redirect(http.StatusSeeOther, c.Request.URL.Path)
	return nil
}

// APIS

func Preview(c *gin.Context) error {
	backend := injector.Backends(c)
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}
	renderedData, err := backend.Render(data)
	if err != nil {
		return err
	}
	c.Data(http.StatusOK, "text/html", renderedData)
	return nil
}
