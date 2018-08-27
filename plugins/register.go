package plugins

import "github.com/jinzhu/gorm"

type Plugin interface {
}

var plugins = map[string]PluginInit{}

type PluginInit func(db *gorm.DB, opts map[string]string) Plugin

func Register(name string, initFunc PluginInit) {
	plugins[name] = initFunc
}

func New(name string, db *gorm.DB, opts map[string]string) (Plugin, error) {
	return plugins[name](db, opts), nil
}
