package modified

import (
	"context"
	"time"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
	"github.com/bluemir/wikinote/internal/pubsub"
	"github.com/pkg/errors"
)

const (
	KeyLastModified = "wikinote.bluemir.me/last-modified"
)

type Options struct {
	Format string `yaml:"format"`
}
type Core struct {
	*Options
	store metadata.IStore
}

var defaultConfig = `
# Last Modified
format: RFC3339
`

func init() {
	plugins.Register("last-modified", New, defaultConfig, &Options{})
}

func New(ctx context.Context, conf any, store metadata.IStore, hub *pubsub.Hub) (plugins.Plugin, error) {
	opt, ok := conf.(*Options)
	if !ok {
		return nil, errors.Errorf("option type not matched: %T", conf)
	}
	return &Core{opt, store}, nil
}

func (c *Core) FileWriteHook(path string, data []byte) ([]byte, error) {
	ctx := context.Background()
	if err := c.store.Save(ctx, path, KeyLastModified, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return data, err
	}
	return data, nil
}
func (c *Core) Footer(path string) ([]byte, error) {
	ctx := context.Background()
	value, err := c.store.Take(ctx, path, KeyLastModified)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return []byte{}, err
	}

	return []byte("last update: " + t.Local().Format(time.RFC3339)), nil
}
func (c *Core) SetConfig(ctx context.Context, conf any) error {
	return nil
}
