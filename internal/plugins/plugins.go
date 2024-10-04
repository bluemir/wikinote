package plugins

import (
	"context"
	"net/http"
)

type Plugin interface {
	SetConfig(ctx context.Context, conf any) error
	//GetConfig(ctx context.Context) (string, error)
}

type PluginReadHook interface {
	FileReadHook(path string, data []byte) ([]byte, error)
}
type PluginFooter interface {
	Footer(path string) ([]byte, error)
}
type PluginHTTPHandler = http.Handler
