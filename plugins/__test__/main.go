package test_plugins

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/plugins"
)

func init() {
	plugins.Register("__test__", New)
}

func New(db *gorm.DB, opts map[string]string) plugins.Plugin {
	logrus.Debugf("test config: %v", opts)
	return &TestPlugin{}
}

type TestPlugin struct {
}
