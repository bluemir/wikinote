package backend

import (
	"io/ioutil"
	"path/filepath"

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
  - verbs:
    - read
    resource:
      kind: "%image,wiki"
`

func loadConfigFile(wikipath string) (*Config, error) {
	configPath := filepath.Join(wikipath, ".app/config.yaml")
	conf := Config{}
	if err := yaml.Unmarshal([]byte(defaultConfig), &conf); err != nil {
		return nil, err
	}

	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		logrus.Warn("config.yaml not exist.", err)
		return &conf, nil
	}

	if err = yaml.Unmarshal(buf, &conf); err != nil {
		logrus.Warn("config.yaml not parsed.", err)
	}

	return &conf, nil
}
