package authz

import (
	yaml "gopkg.in/yaml.v3"

	"github.com/bluemir/wikinote/plugins"
)

func init() {
	plugins.Register("authz", New)
}

func New(core plugins.Core, conf []byte) (plugins.Plugin, error) {
	opt := &Option{}
	err := yaml.Unmarshal(conf, opt)
	if err != nil {
		return nil, err
	}
	return &Authz{core, opt}, nil
}

type Option struct {
	Rules []Rule `yaml:"rules"`
}

type Authz struct {
	plugins.Core
	opts *Option
}

func (authz *Authz) AuthCheck(c *plugins.AuthContext) (plugins.AuthResult, error) {
	for _, rule := range authz.opts.Rules {
		if rule.match(c) {
			return plugins.Accept, nil
		}
	}
	return plugins.Reject, nil
}
