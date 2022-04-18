package errors

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type HTTPError struct {
	Message string `example:"error message"`
}
type HTTPBadRequestError struct {
	HTTPError
	Fields map[string]string
}
type HTTPInternalServerError struct {
	HTTPError
	Errors []string
}

func Json(c *gin.Context) {
	c.Next()

	if len(c.Errors) == 0 {
		return // skip error handler
	}

	// log all errors
	logrus.Trace(c.Errors.String())

	if c.Writer.Written() {
		logrus.Error("response already written. but has error", c.Errors.String())
		return
	}

	// most root cause of error will attach right before abort/return.
	// backword search for error code.
	for i := 0; i < len(c.Errors); i++ {

		err := c.Errors[len(c.Errors)-(i+1)]

		if code, msg := getCodeMessage(err.Err); code > 0 {
			c.JSON(code, msg)
			return
		}
	}

	// final error handler for safe.
	if len(c.Errors) > 0 && !c.Writer.Written() {
		c.JSON(http.StatusInternalServerError, c.Errors.JSON())
	}
}
func getCodeMessage(err error) (int, interface{}) {
	switch v := err.(type) {
	case nil:
		return 0, "unknown error"
	case validator.FieldError:
		return http.StatusNotImplemented, v.Error()
	case validator.ValidationErrors:
		msgs := map[string]string{}
		for _, e := range v {
			msgs[e.Field()] = validationErrorToText(e)
		}
		return http.StatusBadRequest, v.Error()
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

func validationErrorToText(e validator.FieldError) string {

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s", e.Field(), e.Param())
	case "min":
		return fmt.Sprintf("%s must be longer than %s", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("Invalid email format")
	case "len":
		return fmt.Sprintf("%s must be %s characters long", e.Field(), e.Param())
	}
	return fmt.Sprintf("%s is not valid", e.Field())
}

/*

error struct

---
# bind fail
msg: "bad request"
fields:
  aa: aa is required
  bb: bb too short
# if there is header for bot
_compute:
- field: aa
  tag: required
- field: bb
  tag: min

---
# not found
msg: "xxx not found"
---
# internal error
msg: "internal server error"
errors:
- "aaa"
- "bbb"
---
# common
msg: "aaaaa"
*/
