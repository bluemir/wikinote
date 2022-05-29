package errors

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

type HTTPErrorHandlerOption func(int, *HTTPErrorResponse, reqtype.ReqType) (int, *HTTPErrorResponse, reqtype.ReqType)

func WithType(reqType reqtype.ReqType) HTTPErrorHandlerOption {
	return func(code int, res *HTTPErrorResponse, t reqtype.ReqType) (int, *HTTPErrorResponse, reqtype.ReqType) {
		return code, res, reqType
	}
}

func HTTPErrorHandler(c *gin.Context, err error, opts ...HTTPErrorHandlerOption) {
	c.Abort()

	if c.Writer.Written() && c.Writer.Size() > 0 {
		logrus.Trace("response already written")
		return // skip. already written
	}

	code := findErrorCode(err)
	res := makeErrorResponse(err)
	t := reqtype.FindRequestType(c)

	for _, f := range opts {
		code, res, t = f(code, res, t)
	}

	logrus.Debug(code, res)

	if code == http.StatusUnauthorized {
		c.Header(auth.HeaderWWWAuthenticate, "basic realm="+c.Request.URL.Host)
	}

	switch t {
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
