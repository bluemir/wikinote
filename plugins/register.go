package plugins

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Plugin interface {
}

var plugins = map[string]PluginInit{}

type PluginInit func(core Core, conf []byte) (Plugin, error)

func Register(name string, initFunc PluginInit) {
	plugins[name] = initFunc
}
func New(name string, core Core, config []byte) (Plugin, error) {
	log := logrus.WithField("method", "plugin.New")
	log.Tracef("name: %s, opts: %s", name, config)
	if init, ok := plugins[name]; ok {
		return init(core, config)
	}
	log.Tracef("plugins: %#v", plugins)
	return nil, errors.Errorf("Plugin not Found: %s", name)
}
