package plugins

import (
	"html/template"

	"github.com/bluemir/wikinote/pkgs/fileattr"
	"github.com/gin-gonic/gin"
)

type FooterPlugin interface {
	Footer(path string, attr FileAttr) (template.HTML, error)
}
type PostSavePlugin interface {
	OnPostSave(path string, data []byte, attr FileAttr) error
}
type AllHookPlugin interface {
	PostSavePlugin
	//PreSavePlugin
}
type RegisterRouterPlugin interface {
	RegisterRouter(r gin.IRouter)
}
type ContentsPermissionPlugin interface {
	CanView(path string, user string) (bool, error)
	CanEdit(path string, user string) (bool, error)
}
type FileAttr interface {
	Get(key string) (string, error)
	Set(key, value string) error
	All(namespace string) (map[string]string, error)
}
type FileAttrStore = fileattr.Store
type FindOptions = fileattr.Options

const (
	KEY   = fileattr.OrderTypeKey
	VALUE = fileattr.OrderTypeValue

	ASC  = fileattr.OrderDirectionAsc
	DESC = fileattr.OrderDirectionDesc
)
