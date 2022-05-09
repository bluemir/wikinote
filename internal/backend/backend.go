package backend

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/backend/attr"
	"github.com/bluemir/wikinote/internal/backend/files"
	"github.com/bluemir/wikinote/internal/plugins"
)

type Config struct {
	Salt      string                 `yaml:"salt"`
	FrontPage string                 `yaml:"front-page"`
	Plugins   []plugins.PluginConfig `yaml:"plugins"`
	Roles     []auth.Role            `yaml:"roles"`
}
type Backend struct {
	Config   *Config
	db       *gorm.DB
	Auth     *auth.Manager
	Plugin   *plugins.Manager
	FileAttr *attr.Store
	files    *files.FileStore
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

	fileAttr, err := initFileAttr(db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init file attribute module")
	}

	auth, err := initAuth(db, conf.Salt, conf.Roles)
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
	plugin, err := initPlugins(conf.Plugins, fileAttr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init admin user")
	}

	backend := &Backend{
		Config:   conf,
		db:       db,
		FileAttr: fileAttr,
		Auth:     auth,
		Plugin:   plugin,
		files:    store,
	}
	logrus.Trace("backend initailized")

	return backend, nil
}
