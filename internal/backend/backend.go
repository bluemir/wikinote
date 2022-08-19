package backend

import (
	"path/filepath"

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
	Salt        string                 `yaml:"salt"`
	Plugins     []plugins.PluginConfig `yaml:"plugins"`
	Roles       []auth.Role            `yaml:"roles"`
	DefaultRole string                 `yaml:"default-role"`
}
type Backend struct {
	wikipath string
	Config   *Config
	Auth     *auth.Manager

	db     *gorm.DB
	files  *files.FileStore
	attr   *attr.Store
	Plugin *plugins.Manager
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

	attr, err := initFileAttr(db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init file attribute module")
	}

	auth, err := initAuth(db, conf.Salt, conf.Roles, conf.DefaultRole)
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
	plugin, err := initPlugins(conf.Plugins, attr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init plugins")
	}

	backend := &Backend{
		wikipath: wikipath,
		Config:   conf,
		db:       db,
		attr:     attr,
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
