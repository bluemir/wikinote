package backend

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/backend/events"
	"github.com/bluemir/wikinote/internal/backend/files"
	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
	"github.com/bluemir/wikinote/internal/pubsub"
)

type Config struct {
	Salt     string                 `yaml:"salt"`
	Plugins  []plugins.PluginConfig `yaml:"plugins"`
	Metadata metadata.Config        `yaml:"metadata"`
}
type Backend struct {
	wikipath string

	//db  *gorm.DB
	hub *pubsub.Hub

	Config *Config

	Auth     *auth.Manager
	files    *files.FileStore
	Metadata metadata.Store
	events   *events.EventRecoder
	Plugin   *plugins.Manager
}

func New(ctx context.Context, wikipath string, volatileDatabase bool) (*Backend, error) {
	if err := os.MkdirAll(filepath.Join(wikipath, ".app"), 0755); err != nil {
		return nil, errors.WithStack(err)
	}

	dbPath := filepath.Join(wikipath, ".app/wikinote.db")
	if volatileDatabase {
		dbPath = ":memory:"
	}

	db, err := initDB(dbPath)
	if err != nil {
		return nil, err
	}

	hub, err := pubsub.NewHub(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init message hub")
	}

	recoder, err := events.New(ctx, db, hub)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init event recoder")
	}

	auth, err := auth.New(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init auth module")
	}

	mdstore, err := metadata.New(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init metadata module")
	}

	store, err := initFileStore(wikipath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init file store")
	}

	conf := []plugins.PluginConfig{}
	// TODO Load plugin configs

	plugin, err := plugins.New(conf, mdstore)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init plugins")
	}

	logrus.Trace("backend initailized")

	defer hub.Publish("system", Message{"server started"})

	return &Backend{
		wikipath: wikipath,

		//db:  db,
		hub: hub,

		Auth:     auth,
		files:    store,
		Metadata: mdstore,
		events:   recoder,
		Plugin:   plugin,
	}, nil
}

type Message struct {
	Text string
}

func (b *Backend) ConfigPath(path string) string {
	return filepath.Join(b.wikipath, ".app", path)
}
