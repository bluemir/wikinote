package test_plugins

import (
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/plugins"
)

func init() {
	plugins.Register("__test__", New)
}

func New(opts map[string]string, fileAttrStore plugins.FileAttrStore) plugins.Plugin {
	logrus.Debugf("test config: %v", opts)
	return &TestPlugin{}
}

type TestPlugin struct {
}
