package backend

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/backend/files"
	"github.com/bluemir/wikinote/internal/fileattr"
	"github.com/bluemir/wikinote/internal/plugins"
)

func InitConfig() Config {
	return Config{
		AdminUsers: map[string]string{},
	}
}

type Config struct {
	Wikipath   string
	ConfigFile string

	AdminUsers map[string]string

	File struct {
		FrontPage string                 `yaml:"front-page"`
		Plugins   []plugins.PluginConfig `yaml:"plugins"`
		Roles     []auth.Role            `yaml:"roles"`
	}
}
type Backend struct {
	Config   *Config
	Auth     *auth.Manager
	Plugin   *plugins.Manager
	FileAttr *fileattr.Store
	db       *gorm.DB

	files files.FileStore
}

func New(conf *Config) (*Backend, error) {
	// Load config file
	if err := loadConfigFile(conf); err != nil {
		return nil, err
	}

	buf, _ := yaml.Marshal(conf)
	logrus.Debugf("config:\n%s", buf)

	backend := &Backend{Config: conf}

	if err := backend.initDB(); err != nil {
		return nil, err
	}
	if err := backend.initFileAttr(); err != nil {
		return nil, errors.Wrap(err, "failed to init file attribute module")
	}
	if err := backend.initAuth(); err != nil {
		return nil, errors.Wrap(err, "failed to init auth module")
	}
	if err := backend.initAdminUser(); err != nil {
		return nil, errors.Wrap(err, "failed to init admin user")
	}
	if err := backend.initPlugins(); err != nil {
		return nil, errors.Wrap(err, "failed to init admin user")
	}
	if err := backend.initFileStore(); err != nil {
		return nil, errors.Wrap(err, "failed to init file store")
	}

	logrus.Trace("backend initailized")

	return backend, nil
}
