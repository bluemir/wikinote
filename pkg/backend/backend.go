package backend

import (
	"os"
	"path/filepath"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkg/auth"
)

type Config struct {
	Wikipath   string
	ConfigFile string
	RoleFile   string

	File struct {
		FrontPage string `yaml:"front-page"`
		Plugins   []struct {
			Name    string      `yaml:"name"`
			Options interface{} `yaml:"options"`
		} `yaml:"plugins"`
	}
}
type Backend struct {
	Config *Config
	Auth   *auth.Manager
	db     *gorm.DB
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
	if err := db.AutoMigrate(
		&FileAttr{},
	).Error; err != nil {
		return nil, errors.Wrap(err, "auto migrate is failed")
	}

	// Init Auth Module
	authManager, err := auth.New(db, conf.RoleFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init auth module")
	}

	log.Trace("backend initailized")

	return &Backend{conf, authManager, db}, nil
}
