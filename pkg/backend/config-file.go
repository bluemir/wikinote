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
`

func loadConfigFile(conf *Config) error {
	log := logrus.WithField("method", "backend.loadConfigFile")
	if err := yaml.Unmarshal([]byte(defaultConfig), &conf.File); err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(conf.ConfigFile)
	if err != nil {
		log.Warn("config.yaml not exist.", err)
	} else {
		if err = yaml.Unmarshal(buf, &conf.File); err != nil {
			log.Warn("config.yaml not parsed.", err)
		}
	}
	return nil
}
