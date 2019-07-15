package plugins

import (
	"github.com/bluemir/go-utils/auth"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

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
	log := logrus.WithField("method", "plugin.New")
	log.Tracef("name: %s, opts: %#v", name, opts)
	if plugin, ok := plugins[name]; ok {
		return plugin(opts, fileAttrStore, authManager), nil
	}
	log.Tracef("plugins: %#v", plugins)
	return nil, errors.Errorf("Plugin not Found: %s", name)
}
