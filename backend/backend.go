package backend

import (
	"os"
	"time"

	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	"github.com/docker/libkv/store/boltdb"
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

	Store() store.Store

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

	// second init kv store for user
	kv, err := libkv.NewStore(
		store.BOLTDB,
		[]string{wikipath + "/.app/wikinote.kvstore.db"},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
			Bucket:            "wikinote",
		},
	)
	if err != nil {
		return nil, err
	}

	// make backend structor
	b := &backend{
		basePath:   wikipath,
		conf:       conf,
		configPath: configFile,
		kv:         kv,
	}

	// initialize components
	if b.authManager, err = auth.NewManager(); err != nil {
		return nil, err
	}
	if b.fileManager, err = file.New(wikipath); err != nil {
		return nil, err
	}
	if b.userManager, err = user.NewManager(kv, conf); err != nil {
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

	authManager   auth.Manager
	fileManager   file.Manager
	userManager   user.Manager
	pluginManager plugin.Manager
}

func (b *backend) Close() {
	b.Store().Close()
	b.SaveConfig(b.Config())
}

func (b *backend) Config() *config.Config {
	return b.conf
}
func (b *backend) SaveConfig(conf *config.Config) error {
	return b.conf.Save(b.configPath)
}
func (b *backend) Store() store.Store {
	return b.kv
}

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
