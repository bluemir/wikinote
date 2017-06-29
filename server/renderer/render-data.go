package renderer

import (
	"path"
	"strings"

	//"github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

type RenderData struct {
	Path       string
	Title      string
	UserName   string
	UserData   map[string]interface{}
	Breadcrumb []Breadcrumb
	MsgInfo    []interface{}
	MsgWarn    []interface{}
}
type Breadcrumb struct {
	Name string
	Path string
}

func (rd *RenderData) IsLogin() bool {
	return rd.UserName != ""
}

func getWikipath(c *gin.Context) string {
	path := c.Request.URL.Path
	if strings.HasPrefix(path, "/!/") {
		return c.Param("path")
	}
	return path
}
func parseBreadcrumb(p string) []Breadcrumb {
	result := []Breadcrumb{}
	arr := strings.Split(p, "/")
	for index, name := range arr {
		if name == "" {
			continue
		}
		p := path.Join(arr[:index+1]...)
		result = append(result, Breadcrumb{
			Name: name,
			Path: path.Join("/", p),
		})
	}
	return result
}
