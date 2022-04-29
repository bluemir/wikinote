package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/server/middleware/reqtype"
)

func ErrorHandler(c *gin.Context, err error) {
	c.Abort()

	if c.Writer.Written() && c.Writer.Size() > 0 {
		logrus.Trace("response already written")
		return // skip. already written
	}

	code := findErrorCode(err)
	res := makeErrorResponse(err)

	logrus.Debug(code, res)

	if code == http.StatusUnauthorized {
		c.Header(auth.HeaderWWWAuthenticate, "basic realm="+c.Request.URL.Host)
	}

	switch reqtype.FindRequestType(c) {
	case reqtype.API:
		c.JSON(code, res)
		return
	case reqtype.HTML:
		c.HTML(code, getErrorHTMLName(code), res)
		return
	case reqtype.Unknown:
		c.String(code, res.String())
	default:
		c.HTML(code, getErrorHTMLName(code), res)
		return
	}
}
func findErrorCode(err error) int {
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
func getErrorHTMLName(code int) string {
	switch code {
	case http.StatusUnauthorized:
		return "/errors/unauthorized.html"
	case http.StatusNotFound:
		return "/errors/not-found.html"
	default:
		return "/errors/internal-sever-error.html"
	}
}

type HTTPErrorResponse struct {
	Message string   `json:"message"`
	Cause   []string `json:"cause,omitempty"`
}

func makeErrorResponse(err error) *HTTPErrorResponse {
	return &HTTPErrorResponse{
		Message: err.Error(),
	}
}
func (e HTTPErrorResponse) String() string {
	if len(e.Cause) > 0 {
		return e.Message + "\n" + strings.Join(e.Cause, "\n")
	}
	return e.Message
}
