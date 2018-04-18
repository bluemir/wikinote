package renderer

import (
	"path"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/backend/user"
)

// using outside of renderer
type Data map[string]interface{}

// using inside of renderer
type renderData struct {
	context  *gin.Context
	UserData map[string]interface{}

	MsgInfo []interface{}
	MsgWarn []interface{}
}

// using in template
type Breadcrumb struct {
	Name string
	Path string
}

func (d Data) With(c *gin.Context) *renderData {
	return &renderData{
		context:  c,
		UserData: d,
	}
}

func (rd *renderData) WikiPath() string {
	path := rd.context.Request.URL.Path
	if strings.HasPrefix(path, "/!/") {
		return rd.context.Param("path")
	}
	return path
}
func (rd *renderData) Path() string {
	return rd.WikiPath()
}
func (rd *renderData) Breadcrumb() []Breadcrumb {
	return parseBreadcrumb(rd.WikiPath())
}
func (rd *renderData) Username() string {
	u, ok := rd.context.Get("user")
	if !ok {
		return ""
	} else {
		return u.(*user.User).Name
	}
}
func (rd *renderData) IsLogin() bool {
	_, ok := rd.context.Get("user")
	return ok
}
func (rd *renderData) IsSpecial() bool {
	_, isSpecial := rd.context.Get("special-route")
	return isSpecial
}

// prepare data before write body
func (rd *renderData) Prepare() error {
	// TODO check error? validation?

	session := sessions.Default(rd.context)

	rd.MsgInfo = session.Flashes(MSG_INFO)
	rd.MsgWarn = session.Flashes(MSG_WARN)
	session.Save()

	return nil
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
