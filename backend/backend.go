package backend

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/backend/config"
	"github.com/bluemir/wikinote/plugins"
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

		plugins: struct {
			footer        []plugins.FooterPlugin
			afterWikiSave []plugins.AfterWikiSavePlugin
		}{
			footer:        []plugins.FooterPlugin{},
			afterWikiSave: []plugins.AfterWikiSavePlugin{},
		},
	}

	// initialize components
	dbInit(db)
	b.loadPlugins()
	/*if b.authManager, err = auth.NewManager(db); err != nil {
		return nil, err
	}
	if b.fileManager, err = file.New(wikipath, db); err != nil {
		return nil, err
	}
	if b.userManager, err = user.NewManager(db, conf, wikipath); err != nil {
		return nil, err
	}
	if b.pluginManager, err = plugin.New(db, conf.Plugins); err != nil {
		return nil, err
	}
	*/
	logrus.Info("Backend Initialized")

	return b, nil
}

type backend struct {
	basePath   string
	conf       *config.Config
	configPath string
	db         *gorm.DB

	plugins struct {
		footer        []plugins.FooterPlugin
		afterWikiSave []plugins.AfterWikiSavePlugin
	}
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
