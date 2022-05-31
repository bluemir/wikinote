package footer

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/internal/plugins"
)

type Options struct {
	Text string
}
type Core struct {
	*Options
}

func init() {
	plugins.Register("footer", New, &Options{})
}

func New(o interface{}, store *plugins.Store) (plugins.Plugin, error) {
	opt, ok := o.(*Options)
	if !ok {
		return nil, errors.Errorf("option not matched")
	}

	return &Core{opt}, nil
}
func (core *Core) Footer(path string) ([]byte, error) {
	logrus.Info(path, core.Options.Text)

	return []byte(core.Options.Text), nil
}
