package plugins

import (
	"html/template"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/pkgs/auth"
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
type ReadWikiPlugin interface {
	//OnRead(path string, data []byte, attr FileAttr) ([]byte, error)
	OnReadWiki(authCtx *AuthContext, path string, data []byte) ([]byte, error)
}
type AuthzPlugin interface {
	AuthCheck(c *AuthContext) (auth.Result, error)
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

type Core interface {
	// File().Attr().Where().SortBy().Limit().Find()
	// File().Attr().SortBy().Find()
	File() CoreFile
	Auth() CoreAuth
}
type CoreFile interface {
	Attr() CoreFileAttr
}
type CoreFileAttr = fileattr.Store

type CoreAuth interface {
}

type AuthContext = auth.Context
