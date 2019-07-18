package config

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var defaultConfig = `
front-page: front-page.md
auto-backup: false
user:
  default:
    role: viewer
plugins:
  - name: __test__
    options:
      testmap:
        v1: "foo"
        v2: "bar"
`

type Config struct {
	FrontPage  string `yaml:"front-page"`
	AutoBackup bool   `yaml:"auto-backup"`
	User       struct {
		Default struct {
			Role string
		}
	}
	Plugins []struct {
		Name    string      `yaml:"name"`
		Options interface{} `yaml:"options"`
	} `yaml:"plugins"`
}

func ParseConfig(path string) (*Config, error) {
	conf := &Config{}
	defer logrus.Debugf("default: %#v", conf)

	if err := yaml.Unmarshal([]byte(defaultConfig), conf); err != nil {
		return conf, err
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Warn("config.yaml not exist.")
		return conf, nil
	}

	if err = yaml.Unmarshal(buf, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
func (conf *Config) Save(path string) error {
	data, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
