package test

import (
	"context"

	"github.com/pkg/errors"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
	"github.com/bluemir/wikinote/internal/pubsub"
)

type Options struct {
	Foo string
}
type Core struct {
	*Options
}

func init() {
	plugins.Register("_sample", New, defaultConf, &Options{})
}

func New(ctx context.Context, conf any, store metadata.IStore, hub *pubsub.Hub) (plugins.Plugin, error) {
	opt := &Options{}

	opt, ok := conf.(*Options)
	if !ok {
		return nil, errors.Errorf("option type not matched: %T")
	}

	return &Core{opt}, nil
}

// TODO make config store?
var defaultConf = `
# this is sample plugin for test
# DO NOT enable this plugin for production use
`

func (*Core) SetConfig(ctx context.Context, conf any) error {
	_, ok := conf.(*Options)
	if !ok {
		return errors.Errorf("optiontype not matched: %T", conf)
	}
	return nil
}
