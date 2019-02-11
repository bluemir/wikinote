package plugins

import (
	"github.com/bluemir/go-utils/auth"

	"github.com/bluemir/wikinote/pkgs/fileattr"
)

type Plugin interface {
}

var plugins = map[string]PluginInit{}

type PluginInit func(opts map[string]string, store FileAttrStore, auth AuthManager) Plugin

func Register(name string, initFunc PluginInit) {
	plugins[name] = initFunc
}

func New(name string, opts map[string]string, fileAttrStore fileattr.Store, authManager auth.Manager) (Plugin, error) {
	return plugins[name](opts, fileAttrStore, authManager), nil
}
