package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HTTPError struct {
	Message string `example:"error message"`
}
type ErrHandlerOpt func(code int, message string) (int, string)

func APIErrorHandler(c *gin.Context, err error, opts ...ErrHandlerOpt) {
	logrus.Warn(err)

	code, message := defaultCodeMessage(err)

	for _, f := range opts {
		code, message = f(code, message)
	}

	c.JSON(code, HTTPError{Message: message})
	c.Abort()
}
func HTMLErrorHandler(c *gin.Context, err error, opts ...ErrHandlerOpt) {
	logrus.Warn(err)

	code, message := defaultCodeMessage(err)

	for _, f := range opts {
		code, message = f(code, message)
	}

	c.HTML(code, errPage(code), HTTPError{Message: message})
	c.Abort()
}
func errPage(code int) string {
	switch code {
	case http.StatusNotFound:
		return "/errors/not-found.html"
	case http.StatusForbidden:
		return "/errors/forbidden.html"
	case http.StatusUnauthorized:
		return "/errors/unauthorized.html"
	default:
		return "/errors/internal-server-error.html"
	}
}
func defaultCodeMessage(err error) (int, string) {
	switch v := err.(type) {
	case nil:
		return http.StatusInternalServerError, "unknown error"
	default:
		switch v {
		case os.ErrNotExist:
			return http.StatusNotFound, v.Error()
		case gorm.ErrRecordNotFound:
			return http.StatusNotFound, v.Error()
		case gorm.ErrRegistered:
			return http.StatusConflict, v.Error()
		default:
			return http.StatusInternalServerError, v.Error()
		}
	}
}
func withCode(code int) ErrHandlerOpt {
	return func(c int, message string) (int, string) {
		return code, message
	}
}
func withMessage(format string, args ...interface{}) ErrHandlerOpt {
	return func(code int, m string) (int, string) {
		return code, fmt.Sprintf(format, args...)
	}
}
func withAdditionalMessage(format string, args ...interface{}) ErrHandlerOpt {
	message := fmt.Sprintf(format, args...)
	return func(code int, m string) (int, string) {
		return code, fmt.Sprintf("%s: %s", message, m)
	}
}
