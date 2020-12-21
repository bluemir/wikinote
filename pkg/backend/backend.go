package backend

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/bluemir/wikinote/pkg/auth"
	"github.com/bluemir/wikinote/pkg/fileattr"
	"github.com/bluemir/wikinote/pkg/plugins"
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
}

func New(conf *Config) (*Backend, error) {
	// Load config file
	if err := loadConfigFile(conf); err != nil {
		return nil, err
	}

	buf, _ := yaml.Marshal(conf)
	logrus.Debugf("config:\n%s", buf)

	// Init DB
	dbPath := filepath.Join(conf.Wikipath, ".app/wikinote.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, err
	}
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}
	db.DB().SetMaxOpenConns(1)

	fa, err := fileattr.New(db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init file attribute module")
	}

	// Init Auth Module
	authManager, err := auth.New(db, conf.File.Roles)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init auth module")
	}

	for name, key := range conf.AdminUsers {
		if key == "" {
			key = xid.New().String()
			logrus.Warnf("generate key: '%s' '%s'", name, key)
		}
		if err := authManager.EnsureUser(name, map[string]string{
			"role/root": "true",
		}); err != nil {
			return nil, err
		}
		if err := authManager.RevokeTokenAll(name); err != nil {
			return nil, err
		}
		if _, err := authManager.IssueToken(name, key); err != nil {
			return nil, err
		}
	}

	// Init Plugins
	pluginManager, err := plugins.New(conf.File.Plugins, fa)
	if err != nil {
		return nil, err
	}

	logrus.Trace("backend initailized")

	return &Backend{
		Config:   conf,
		Auth:     authManager,
		Plugin:   pluginManager,
		FileAttr: fa,
		db:       db,
	}, nil
}
