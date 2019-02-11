package plugins

import (
	"html/template"

	"github.com/bluemir/go-utils/auth"
	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/pkgs/fileattr"
)

// Plugins
type FooterPlugin interface {
	Footer(path string, attr FileAttr) (template.HTML, error)
}
type PreSavePlugin interface {
	OnPreSave(path string, data []byte, attr FileAttr) ([]byte, error)
}
type PostSavePlugin interface {
	OnPostSave(path string, data []byte, attr FileAttr) error
}
type ReadPlugin interface {
	OnRead(path string, data []byte, attr FileAttr) ([]byte, error)
}
type FilePermissionPlugin interface {
	TryRead(path string, user interface{}, attr FileAttr) error
	TryWrite(path string, user interface{}, attr FileAttr) error
}
type RegisterRouterPlugin interface {
	RegisterRouter(r gin.IRouter)
}

// File Attr

type FileAttr = fileattr.PathClause
type FileAttrStore = fileattr.Store
type FindOptions = fileattr.Options

const (
	KEY   = fileattr.OrderTypeKey
	VALUE = fileattr.OrderTypeValue

	ASC  = fileattr.OrderDirectionAsc
	DESC = fileattr.OrderDirectionDesc
)

// Auth

type AuthManager auth.Manager

// context
type Components interface {
	Auth() AuthManager
	AttrStore() FileAttrStore
}
