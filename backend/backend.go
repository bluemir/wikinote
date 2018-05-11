package backend

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend/config"
)

type Backend interface {
	Config() *config.Config
	SaveConfig(conf *config.Config) error

	File() FileClause
	User() UserClause
	Auth() AuthClause
	Plugin() PluginClause

	// renderer
	Render(data []byte) ([]byte, error)
	//Plugins().Footer() ([]RenderResult)

	Close()
}

type Options struct {
	Wikipath   string
	ConfigFile string
	Version    string
}

func New(o *Options) (Backend, error) {
	logrus.Infof("VERSION: %s", o.Version)
	// first, parse config from file
	wikipath := os.ExpandEnv(o.Wikipath)
	configFile := os.ExpandEnv(o.ConfigFile)
	conf, err := config.ParseConfig(configFile)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(wikipath+"/.app", 0755)
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open("sqlite3", wikipath+"/.app/wikinote.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	// make backend structor
	b := &backend{
		basePath: wikipath,
		conf:     conf,
		db:       db,
		plugins:  nil,
	}

	// initialize components
	err = dbInit(db)
	if err != nil {
		return nil, err
	}
	pl, err := loadPlugins(db, conf)
	if err != nil {
		return nil, err
	}
	b.plugins = pl

	logrus.Info("Backend Initialized")

	return b, nil
}

type backend struct {
	basePath   string
	conf       *config.Config
	configPath string
	db         *gorm.DB

	plugins *pluginList
}

func (b *backend) Close() {
	b.db.Close()
	b.SaveConfig(b.Config())
}

func (b *backend) Config() *config.Config {
	return b.conf
}
func (b *backend) SaveConfig(conf *config.Config) error {
	return b.conf.Save(b.configPath)
}
func (b *backend) Auth() AuthClause {
	return (*authClause)(b)
}
func (b *backend) File() FileClause {
	return (*fileClause)(b)
}
func (b *backend) User() UserClause {
	return (*userClause)(b)
}
func (b *backend) Plugin() PluginClause {
	return (*pluginClause)(b)
}
