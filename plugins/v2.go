package plugins

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var pluginsV2 = map[string]PluginInitV2{}

type PluginInitV2 func(core Core, conf []byte) (Plugin, error)

func RegisterV2(name string, initFunc PluginInitV2) {
	pluginsV2[name] = initFunc
}

func NewV2(name string, core Core, config []byte) (Plugin, error) {
	log := logrus.WithField("method", "plugin.NewV2")
	log.Tracef("name: %s, opts: %s", name, config)
	if init, ok := pluginsV2[name]; ok {
		return init(core, config)
	}
	log.Tracef("plugins: %#v", plugins)
	return nil, errors.Errorf("Plugin not Found: %s", name)
}

type Core interface {
	// File().Attr().Where().SortBy().Limit().Find()
	// File().Attr().SortBy().Find()
	File() CoreFile
	Auth() CoreAuth
}
type Config interface {
	//UnmarshalYAML(value *yaml.Node) error
}
type CoreFile interface {
}
type CoreAuth interface {
}

/*
type Core interface {
	File() interface {
		Attr() interface {
			Where(*Options) WhereClause
			SortBy(OrderType, OrderDirection) SortByClause
			Limit(int) LimitClause
			Find() ([]FileAttr, error)
		}
	}
	Auth() interface {
		SetUserAttr()
		SetTokenAttr()
	}
}
*/
