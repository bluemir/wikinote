package authz

import (
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"

	"github.com/bluemir/wikinote/plugins"
)

func init() {
	plugins.RegisterV2("authz", New)
}

func New(core plugins.Core, conf []byte) (plugins.Plugin, error) {
	logrus.Debugf("#############################\n%s", conf)
	opt := &Option{}
	err := yaml.Unmarshal(conf, opt)
	if err != nil {
		return nil, err
	}
	return &Authz{core}, nil
}

type Option struct {
	Rules []struct {
		Object  []KV
		Subject []KV
		Action  []string
	} `yaml:"rules"`
}
type KV struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Authz struct {
	plugins.Core
}
