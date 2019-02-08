package plugins

import (
	"github.com/bluemir/wikinote/pkgs/fileattr"
)

type Plugin interface {
}

var plugins = map[string]PluginInit{}

type PluginInit func(opts map[string]string, store FileAttrStore) Plugin

func Register(name string, initFunc PluginInit) {
	plugins[name] = initFunc
}

func New(name string, opts map[string]string, fileAttrStore fileattr.Store) (Plugin, error) {
	return plugins[name](opts, fileAttrStore), nil
}
