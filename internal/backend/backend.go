package backend

import (
	"context"
	"encoding/gob"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/bluemir/wikinote/internal/auth"
	"github.com/bluemir/wikinote/internal/backend/files"
	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/events"
	"github.com/bluemir/wikinote/internal/plugins"
)

type Config struct {
	Salt     string                 `yaml:"salt"`
	Plugins  []plugins.PluginConfig `yaml:"plugins"`
	Auth     auth.Config            `yaml:"auth"`
	Metadata metadata.Config        `yaml:"metadata"`
}
type Backend struct {
	wikipath string
	Config   *Config

	Auth     *auth.Manager
	Metadata metadata.Store
	Plugin   *plugins.Manager
	files    *files.FileStore
	hub      *events.Hub

	db *gorm.DB
}

func New(ctx context.Context, wikipath string, volatileDatabase bool) (*Backend, error) {
	dbPath := filepath.Join(wikipath, ".app/wikinote.db")
	if volatileDatabase {
		dbPath = ":memory:"
	}

	db, err := initDB(dbPath)
	if err != nil {
		return nil, err
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

	gob.Register(Message{})

	hub, err := events.NewHub(context.TODO(), db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init message hub")
	}

	if err := hub.Fire("user/bluemir", Message{"server started"}); err != nil {
		return nil, err
	}

	backend := &Backend{
		wikipath: wikipath,
		db:       db,
		Metadata: mdstore,
		Auth:     auth,
		Plugin:   plugin,
		files:    store,
		hub:      hub,
	}
	logrus.Trace("backend initailized")

	return backend, nil
}

type Message struct {
	Text string
}

func (b *Backend) ConfigPath(path string) string {
	return filepath.Join(b.wikipath, ".app", path)
}
