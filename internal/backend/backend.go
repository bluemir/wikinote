package backend

import (
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/backend/files"
	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
)

type Config struct {
	Salt     string                 `yaml:"salt"`
	Plugins  []plugins.PluginConfig `yaml:"plugins"`
	Auth     auth.Config            `yaml:"auth"`
	Metadata metadata.Config        `yaml:"metadata"`
}
type Backend struct {
	wikipath string
	Config   *Config

	Auth     *auth.Manager
	Metadata metadata.Store
	Plugin   *plugins.Manager
	files    *files.FileStore

	db *gorm.DB
}

func New(wikipath string, users map[string]string) (*Backend, error) {
	// Load config file
	conf, err := loadConfigFile(wikipath)
	if err != nil {
		return nil, err
	}

	buf, _ := yaml.Marshal(conf)
	logrus.Debugf("config:\n%s", buf)

	db, err := initDB(wikipath)
	if err != nil {
		return nil, err
	}

	if conf.Metadata.File != nil && conf.Metadata.File.Path == "" {
		conf.Metadata.File.Path = wikipath
	}
	if conf.Metadata.Gorm != nil {
		conf.Metadata.Gorm.DB = db
	}
	mdstore, err := metadata.New(&conf.Metadata)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init metadata module")
	}

	auth, err := initAuth(db, conf.Salt, &conf.Auth)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init auth module")
	}

	if err := initAdminUser(auth, users); err != nil {
		return nil, errors.Wrap(err, "failed to init admin user")
	}
	store, err := initFileStore(wikipath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init file store")
	}
	plugin, err := plugins.New(conf.Plugins, mdstore)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init plugins")
	}

	backend := &Backend{
		wikipath: wikipath,
		Config:   conf,
		db:       db,
		Metadata: mdstore,
		Auth:     auth,
		Plugin:   plugin,
		files:    store,
	}
	logrus.Trace("backend initailized")

	return backend, nil
}
func (b *Backend) ConfigPath(path string) string {
	return filepath.Join(b.wikipath, ".app", path)
}
