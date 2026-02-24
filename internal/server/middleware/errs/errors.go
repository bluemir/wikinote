package errs

import (
	"mime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HttpError is an interface for errors that carry an HTTP status code.
type HttpError interface {
	Code() int
}

const (
	formatHTML    = "html"
	formatJSON    = "json"
	formatProblem = "problem+json"
	formatText    = "text"
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

	// Last one is most important
	err := c.Errors.Last().Err

	statusCode := code(err)
	if c.Writer.Written() {
		logrus.Debugf("Response code already written, expected '%d', but it was '%d'", statusCode, c.Writer.Status())
		statusCode = c.Writer.Status()
	}

	logrus.Tracef("%T %+v, %d", err, err, statusCode)

	if statusCode >= 500 {
		logrus.Warnf("Server Error. code: %d, %s", statusCode, err)
	}

	format := negotiateResponseFormat(c)
	render(c, format, statusCode, err)
}

func negotiateResponseFormat(c *gin.Context) string {
	// Check format adjusters first
	for _, adjuster := range formatAdjusters {
		if format := adjuster(c); format != "" {
			return format
		}
	}

	// Content negotiation via Accept header
	for accept := range strings.SplitSeq(c.Request.Header.Get("Accept"), ",") {
		t, _, e := mime.ParseMediaType(accept)
		if e != nil {
			logrus.Error(e)
			continue
		}

		switch t {
		case "application/problem+json":
			return formatProblem
		case "application/json":
			return formatJSON
		case "text/html", "*/*":
			return formatHTML
		case "text/plain":
			return formatText
		}
	}

	return formatText
}
