package plugins

import "github.com/gin-gonic/gin"

type Plugin interface {
}

type PluginReadHook interface {
	FileReadHook(path string, data []byte) ([]byte, error)
}
type PluginWriteHook interface {
	FileWriteHook(path string, data []byte) ([]byte, error)
}
type PluginFooter interface {
	Footer(path string) ([]byte, error)
}
type PluginRoute interface {
	Route(r gin.IRouter) error
}
