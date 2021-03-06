package backend

import (
	"os"

	"github.com/bluemir/go-utils/auth"
	_ "github.com/bluemir/go-utils/auth/gorm"
	"github.com/ghodss/yaml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkgs/config"
	"github.com/bluemir/wikinote/pkgs/fileattr"
)

type Backend interface {
	Config() *config.Config
	SaveConfig(conf *config.Config) error

	File() FileClause
	Auth() auth.Manager
	Plugin() PluginClause

	// renderer
	Render(data []byte) ([]byte, error)

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

	authMng, first, err := auth.New(&auth.Options{
		StoreDriver: "gorm",
		DefaultRole: conf.User.Default.Role,
		RootRole:    "root",
		DriverOpts: map[string]interface{}{
			"db": db,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to init auth manager")
	}

	if first {
		// add default rules...
		rules := &struct {
			Rules map[string][]string
		}{}
		err := yaml.Unmarshal([]byte(defaultRule), rules)
		if err != nil {
			return nil, errors.Wrap(err, "parse default rule failed")
		}

		for k, v := range rules.Rules {
			err := authMng.PutRule(k, v...)
			if err != nil {
				return nil, errors.Wrap(err, "failed to put default rule")
			}
		}
	}

	token, err := authMng.Root("root")
	if err != nil {
		return nil, errors.Wrap(err, "failed to ensure root")
	}

	// QUESTION save file or just print stdout?
	logrus.Infof("Root Token: %s", token)

	// FileAttrStore
	fas, err := fileattr.NewStore(db)
	if err != nil {
		return nil, err
	}

	// make backend structor
	b := &backend{
		basePath:      wikipath,
		conf:          conf,
		db:            db,
		auth:          authMng,
		fileAttrStore: fas,
		plugins:       nil,
	}

	// initialize components
	pl, err := loadPlugins(conf, fas, authMng)
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

	auth          auth.Manager
	fileAttrStore fileattr.Store

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
	return &fileClause{b}
}

func (b *backend) Plugin() PluginClause {
	return &pluginClause{b}
}
