package renderer

import (
	"path"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/pkgs/auth"
)

const (
	KEY_TOKEN         = "token"
	KEY_SPECIAL_ROUTE = "special-route"
)

// using outside of renderer
type Data map[string]interface{}

func (d Data) With(c *gin.Context) *userData {
	return &userData{
		context: c,
		data:    d,
	}
}

type userData struct {
	context *gin.Context
	data    Data
}

// Pack is convert data from context to render data
func (ud *userData) Pack() (*renderData, error) {
	var token *auth.Token
	if t, ok := ud.context.Get("token"); ok {
		token = t.(*auth.Token)
	}

	session := sessions.Default(ud.context)

	msgInfo := session.Flashes(MSG_INFO)
	msgWarn := session.Flashes(MSG_WARN)
	session.Save()

	return &renderData{
		context:  ud.context,
		UserData: ud.data,
		MsgInfo:  msgInfo,
		MsgWarn:  msgWarn,
		Token:    token,
	}, nil
}

// using inside of renderer
type renderData struct {
	context  *gin.Context
	UserData Data
	MsgInfo  []interface{}
	MsgWarn  []interface{}
	Token    *auth.Token
}

// using in template
type Breadcrumb struct {
	Name string
	Path string
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
	p := rd.WikiPath()
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
func (rd *renderData) Username() string {
	if rd.Token == nil {
		return ""
	}
	return rd.Token.UserName
}
func (rd *renderData) IsLogin() bool {
	return rd.Token != nil
}
func (rd *renderData) IsSpecial() bool {
	_, isSpecial := rd.context.Get("special-route")
	return isSpecial
}
