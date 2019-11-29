package backend

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/rs/xid"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkg/auth"
	"github.com/bluemir/wikinote/pkg/fileattr"
	"github.com/bluemir/wikinote/pkg/plugins"
)

type Config struct {
	Wikipath   string
	ConfigFile string
	RoleFile   string
	RootUser   string

	File struct {
		FrontPage string                 `yaml:"front-page"`
		Plugins   []plugins.PluginConfig `yaml:"plugins"`
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
	log := logrus.WithField("method", "backend.New")

	// Load config file
	if err := loadConfigFile(conf); err != nil {
		return nil, err
	}

	log.Debug(conf)

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
	authManager, err := auth.New(db, conf.RoleFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init auth module")
	}

	if conf.RootUser != "" {
		key := xid.New().String()
		if err := authManager.EnsureUser(conf.RootUser, map[string]string{
			"role/root": "true",
		}); err != nil {
			return nil, err
		}
		if err := authManager.RevokeTokenAll(conf.RootUser); err != nil {
			return nil, err
		}
		if _, err := authManager.IssueToken(conf.RootUser, key); err != nil {
			return nil, err
		}
		log.Warnf("root key: '%s'", key)
	}

	// Init Plugins
	pluginManager, err := plugins.New(conf.File.Plugins, fa)
	if err != nil {
		return nil, err
	}

	log.Trace("backend initailized")

	return &Backend{
		Config:   conf,
		Auth:     authManager,
		Plugin:   pluginManager,
		FileAttr: fa,
		db:       db,
	}, nil
}
