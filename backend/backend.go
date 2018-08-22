package backend

import (
	"os"

	"github.com/bluemir/go-utils/auth"
	_ "github.com/bluemir/go-utils/auth/gorm"
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
	// User() UserClause
	Auth() auth.Manager
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
	//db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	authMng, err := auth.New(&auth.Options{
		StoreDriver: "gorm",
		DefaultRole: "editor",
		RootRole:    "root",
		DriverOpts: map[string]interface{}{
			"db": db,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to init auth manager")
	}

	err = UserInit(authMng)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init auth manager")
	}

	// make backend structor
	b := &backend{
		basePath: wikipath,
		conf:     conf,
		db:       db,
		auth:     authMng,
		plugins:  nil,
	}

	// initialize components
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

	auth auth.Manager

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
func (b *backend) Auth() auth.Manager {
	return b.auth
}
func (b *backend) File() FileClause {
	return (*fileClause)(b)
}

func (b *backend) Plugin() PluginClause {
	return (*pluginClause)(b)
}
