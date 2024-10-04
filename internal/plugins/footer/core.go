package footer

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
	"github.com/bluemir/wikinote/internal/pubsub"
)

type Options struct {
	Text string `yaml:"text"`
}
type Core struct {
	*Options
}

var defaultConfig = `
# 
text: 
`

func init() {
	plugins.Register("footer", New, defaultConfig, &Options{})
}

func New(ctx context.Context, conf any, store metadata.IStore, hub *pubsub.Hub) (plugins.Plugin, error) {
	opt := &Options{}
	opt, ok := conf.(*Options)
	if !ok {
		return nil, errors.Errorf("option type not matched: %T", conf)
	}

	return &Core{opt}, nil
}
func (core *Core) SetConfig(ctx context.Context, conf any) error {

	opt, ok := conf.(*Options)
	if !ok {
		return errors.Errorf("option type not matched: %T", conf)
	}

	core.Options = opt

	return nil
}
func (core *Core) Footer(path string) ([]byte, error) {
	logrus.Info(path, core.Options.Text)

	return []byte(core.Options.Text), nil
}
