package backend

import (
	"io/ioutil"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

var defaultConfig = `
front-page: front-page.md
plugins:
- name: __test__
  options:
    testmap:
      v1: "foo"
      v2: "bar"
roles:
- name: default
  rules:
  - actions:
    - read
    objects:
    - 'kind == "wiki"'
    - 'kind == "image"'
`

func loadConfigFile(conf *Config) error {
	if err := yaml.Unmarshal([]byte(defaultConfig), &conf.File); err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(conf.ConfigFile)
	if err != nil {
		logrus.Warn("config.yaml not exist.", err)
	} else {
		if err = yaml.Unmarshal(buf, &conf.File); err != nil {
			logrus.Warn("config.yaml not parsed.", err)
		}
	}
	return nil
}
