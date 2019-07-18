package test_plugins

import (
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/plugins"
)

func init() {
	plugins.Register("__test__", New)
}

func New(core plugins.Core, opts []byte) (plugins.Plugin, error) {
	logrus.Debugf("test config: %v", opts)
	return &TestPlugin{}, nil
}

type TestPlugin struct {
}
