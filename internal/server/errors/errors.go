package errors

import (
	"mime"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/auth"
)

type HTTPErrorHandlerOption func(ctx *HTTPErrorContext)

func WithHeader(key, value string) HTTPErrorHandlerOption {
	return func(ctx *HTTPErrorContext) {
		ctx.Headers.Add(key, value)
	}
}
func WithAuthHeader(ctx *HTTPErrorContext) {
	if errors.Is(ctx.Error, auth.ErrUnauthorized) {
		ctx.Header(auth.HeaderWWWAuthenticate, "basic realm="+ctx.Request.URL.Host)
	}
}

type HTTPErrorContext struct {
	*gin.Context
	Code     int
	Error    error
	Response *HTTPErrorResponse
	Headers  http.Header
}

func HTTPErrorHandler(c *gin.Context, err error, opts ...HTTPErrorHandlerOption) {
	c.Abort()

	if c.Writer.Written() && c.Writer.Size() > 0 {
		logrus.Trace("response already written")
		return // skip. already written
	}

	ctx := &HTTPErrorContext{
		Context:  c,
		Error:    err,
		Code:     findErrorCode(err),
		Response: makeErrorResponse(err),
	}

	for _, f := range opts {
		f(ctx)
	}

	logrus.Debugf("%+v", ctx)

	t, _, err := mime.ParseMediaType(c.Request.Header.Get("Accepted"))
	if err != nil {
		logrus.Error(err)
	}

	switch t {
	case "application/json":
		c.JSON(ctx.Code, ctx.Response)
		return
	default:
		c.HTML(ctx.Code, getErrorHTMLName(ctx.Code), ctx.Response)
		return
	}
}
func findErrorCode(err error) int {
	switch {
	case errors.Is(err, auth.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, gorm.ErrRecordNotFound):
		return http.StatusNotFound
	case err == auth.ErrForbidden:
		return http.StatusForbidden
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
	case http.StatusForbidden:
		return "/errors/forbidden.html"
	default:
		return "/errors/internal-server-error.html"
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
