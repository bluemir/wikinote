package errs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/bluemir/wikinote/internal/server/handler"
)

// ProblemDetail represents an RFC 9457 Problem Details object.
type ProblemDetail struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

func render(c *gin.Context, format string, code int, err error) {
	switch format {
	case formatProblem:
		c.Header("Content-Type", "application/problem+json")
		c.JSON(code, ProblemDetail{
			Type:   "about:blank",
			Title:  http.StatusText(code),
			Status: code,
			Detail: err.Error(),
		})
	case formatJSON:
		c.JSON(code, gin.H{
			"errors": c.Errors,
		})
	case formatHTML:
		c.HTML(code, htmlName(code, err), handler.With(c, handler.KeyValues{
			"message": err.Error(),
			"errors":  c.Errors,
		}))
	default: // formatText
		c.String(code, "%#v", c.Errors)
	}
}

func htmlName(code int, err error) string {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
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
