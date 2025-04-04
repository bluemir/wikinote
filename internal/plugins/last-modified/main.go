package modified

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/backend/events"
	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
	"github.com/bluemir/wikinote/internal/pubsub"
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
`

func init() {
	plugins.Register("last-modified", New, defaultConfig, &Options{})
}

func New(ctx context.Context, conf any, store metadata.IStore, hub *pubsub.Hub) (plugins.Plugin, error) {
	opt, ok := conf.(*Options)
	if !ok {
		return nil, errors.Errorf("option type not matched: %T", conf)
	}

	go handleFileWrite(ctx, store, hub)

	logrus.Trace("last-modified enabled")

	return &Core{opt, store}, nil
}
func handleFileWrite(ctx context.Context, store metadata.IStore, hub *pubsub.Hub) {
	logrus.Tracef("%+v", hub)

	ch := hub.Watch(events.KindFileWritten, ctx.Done())

	logrus.Tracef("%+v", hub)

	for msg := range ch {
		evt := msg.Detail.(events.FileWritten)
		logrus.Trace(evt)

		if err := store.Save(ctx, evt.Path, KeyLastModified, time.Now().UTC().Format(time.RFC3339)); err != nil {
			logrus.Warn(err)
			hub.Publish("error", err)
		}
		logrus.WithField("path", evt.Path).WithField("time", time.Now()).Trace("last-modified updated")
	}
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
