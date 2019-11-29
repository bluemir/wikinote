package recent

import (
	"github.com/bluemir/wikinote/pkg/plugins"
	"github.com/pkg/errors"
)

type Options struct {
}
type Recents struct {
	opt *Options
}

func init() {
	plugins.Register("recent-changes", New, &Options{})
}
func New(o interface{}) (plugins.Plugin, error) {
	opt, ok := o.(*Options)
	if !ok {
		return nil, errors.Errorf("option not matched")
	}

	return &Recents{opt}, nil
}

func (r *Recents) FileWriteHook(path string, data []byte) ([]byte, error) {
	// TODO register lastModify

	return data, nil
}
