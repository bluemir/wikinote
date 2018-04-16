package backend

import (
	"os"

	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend/auth"
	"github.com/bluemir/wikinote/backend/config"
	"github.com/bluemir/wikinote/backend/file"
	"github.com/bluemir/wikinote/backend/plugin"
	"github.com/bluemir/wikinote/backend/user"
)

type Backend interface {
	Config() *config.Config
	SaveConfig(conf *config.Config) error

	File() file.Manager
	User() user.Manager
	Auth() auth.Manager
	Plugin() plugin.Manager

	// renderer
	Render(data []byte) ([]byte, error)

	//Plugins().Footer() ([]RenderResult)

	//Store() store.Store

	Close()
}

type Options struct {
	Wikipath   string
	ConfigFile string
	Version    string
}

func init() {
	boltdb.Register()
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

	if err != nil {
		return nil, err
	}
	db, err := gorm.Open("sqlite3", wikipath+"/.app/wikinote.db")
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	// make backend structor
	b := &backend{
		basePath:   wikipath,
		conf:       conf,
		configPath: configFile,
		db:         db,
	}

	// initialize components
	if b.authManager, err = auth.NewManager(); err != nil {
		return nil, err
	}
	if b.fileManager, err = file.New(wikipath); err != nil {
		return nil, err
	}
	if b.userManager, err = user.NewManager(db, conf); err != nil {
		return nil, err
	}
	if b.pluginManager, err = plugin.New(conf.Plugins); err != nil {
		return nil, err
	}
	logrus.Info("Backend Initialized")

	return b, nil
}

type backend struct {
	basePath   string
	conf       *config.Config
	configPath string
	kv         store.Store
	db         *gorm.DB

	authManager   auth.Manager
	fileManager   file.Manager
	userManager   user.Manager
	pluginManager plugin.Manager
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

/*func (b *backend) Store() store.Store {
	return b.kv
}*/

func (b *backend) Auth() auth.Manager {
	return b.authManager
}
func (b *backend) File() file.Manager {
	return b.fileManager
}
func (b *backend) User() user.Manager {
	return b.userManager
}
func (b *backend) Plugin() plugin.Manager {
	return b.pluginManager
}
