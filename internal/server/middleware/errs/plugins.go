package errs

import "github.com/gin-gonic/gin"

// FormatAdjuster is a function that can override the response format
// based on the request context. Return an empty string to skip.
type FormatAdjuster func(c *gin.Context) string

var formatAdjusters []FormatAdjuster

// RegisterFormatAdjuster adds a format adjuster plugin.
func RegisterFormatAdjuster(adj FormatAdjuster) {
	formatAdjusters = append(formatAdjusters, adj)
}
