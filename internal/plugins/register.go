package plugins

import (
	"reflect"

	"github.com/pkg/errors"
)

type PluginInitDriver struct {
	Init    PluginInit
	Type    reflect.Type
	Default string
}

var drivers = map[string]PluginInitDriver{}

func Register(name string, init PluginInit, defaultConfig string, configType any) {
	drivers[name] = PluginInitDriver{
		Init:    init,
		Type:    reflect.TypeOf(configType),
		Default: defaultConfig,
	}
}

func (d *PluginInitDriver) newConfig() any {
	return reflect.New(d.Type).Interface()
}
func getDriver(name string) (*PluginInitDriver, error) {
	d, ok := drivers[name]
	if !ok {
		return nil, errors.Errorf("plugin not found")
	}
	return &d, nil
}
