package plugintest

import (
	"github.com/pkg/errors"

	"github.com/bluemir/wikinote/internal/plugins"
)

type Options struct {
}
type Core struct {
	*Options
}

func init() {
	plugins.Register("__test__", New, &Options{})
}

func New(o interface{}, store *plugins.Store) (plugins.Plugin, error) {
	opt, ok := o.(*Options)
	if !ok {
		return nil, errors.Errorf("option not matched")
	}

	return &Core{opt}, nil
}
