package errors

import (
	"mime"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/handler"
)

func Middleware(c *gin.Context) {
	c.Next()

	errs := c.Errors.ByType(gin.ErrorTypeAny)
	if len(errs) == 0 {
		return
	}

	if c.Writer.Written() && c.Writer.Size() > 0 {
		logrus.Tracef("response already written: %d, %s", c.Writer.Status(), c.Errors.String())
		return // skip. already written
	}

	logrus.Tracef("%d", c.Writer.Status())

	// Last one is most important
	err := c.Errors.Last().Err

	code := code(err)
	if c.Writer.Written() {
		logrus.Debugf("Response code already written, expected '%d', but it was '%d'", code, c.Writer.Status())
		code = c.Writer.Status()
	}

	logrus.Tracef("%T %+v, %d", err, err, code)

	if code >= 500 {
		logrus.Warnf("Server Error. code: %d, %s", code, err)
	}

	// with header or without header, or other processer/ maybe hook? depend on error type? or just code
	for _, accept := range strings.Split(c.Request.Header.Get("Accept"), ",") {
		t, _, e := mime.ParseMediaType(accept)
		if e != nil {
			logrus.Error(e)
			continue
		}

		switch t {
		case "application/json":
			// TODO make response json
			c.JSON(code, gin.H{
				"errors": c.Errors,
			})
			return
		case "text/html", "*/*":
			/*
				if code == http.StatusUnauthorized {
					c.Header(auth.LoginHeader(c.Request)) // for basic auth
				}
			*/
			c.Writer.Status()
			c.HTML(code, htmlName(code, err), handler.With(c, handler.KeyValues{
				"message": err.Error(),
				"errors":  c.Errors,
			}))
			return
		case "text/plain":
			c.String(code, "%#v", c.Errors)
			return
		}
	}
	c.String(code, "%#v", c.Errors)
}

type HttpError interface {
	Code() int
}

func code(err error) int {
	logrus.Tracef("%T", err)

	// errors.Is check same value, but errors.As check only its type.
	switch {
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return http.StatusConflict
	case errors.Is(err, validator.ValidationErrors{}):
		return http.StatusBadRequest
	case errors.Is(err, auth.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, auth.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, os.ErrNotExist):
		return http.StatusNotFound
	case errors.As(err, &sqlite3.Error{}):
		e := sqlite3.Error{}
		errors.As(err, &e)
		switch e.ExtendedCode {
		case sqlite3.ErrConstraintUnique:
			return http.StatusConflict
		default:
			return http.StatusNotImplemented
		}
	}

	// try to call code function
	if e, ok := err.(HttpError); ok {
		return e.Code()
	}

	// finally check string match
	logrus.Trace(err.Error())
	switch {
	case strings.HasPrefix(err.Error(), "html/template: ") && strings.HasSuffix(err.Error(), " is undefined"):
		//html/template: ".*" is undefined
		return http.StatusNotImplemented
	}

	return http.StatusInternalServerError
}
func htmlName(code int, err error) string {
	switch {
	//override
	case errors.Is(err, validator.ValidationErrors{}):
		return "errors/bad-request.html"
	}

	switch code {
	case http.StatusBadRequest:
		return "errors/bad-request.html"
	case http.StatusUnauthorized:
		return "errors/unauthorized.html"
	case http.StatusForbidden:
		return "errors/forbidden.html"
	case http.StatusNotFound:
		return "errors/not-found.html"
	}

	return "errors/internal-server-error.html"
}
