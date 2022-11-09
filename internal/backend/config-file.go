package backend

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

var defaultConfig = `
auth:
  group:
    unauthoized: guest
    newcomer:
    - viewer
  binding:
    "group/guest":
    - viewer
metadata:
  file:
    ext: .metadata
plugins:
- name: footer
  options:
    text: |
      powered by wikinote
`

func loadConfigFile(wikipath string) (*Config, error) {
	configPath := filepath.Join(wikipath, ".app/config.yaml")
	conf := Config{}

	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		switch {
		case errors.Is(err, os.ErrNotExist):
			logrus.Warn("config.yaml not exist.", err)
			if err := ioutil.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
				return nil, err
			}
			buf = []byte(defaultConfig)
		default:
			return nil, err
		}
	}

	if err = yaml.Unmarshal(buf, &conf); err != nil {
		logrus.Warn("config.yaml is not valid: ", err)
	}

	return &conf, nil
}
