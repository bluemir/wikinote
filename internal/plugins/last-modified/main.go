package modified

import (
	"time"

	"github.com/pkg/errors"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/bluemir/wikinote/internal/plugins"
)

const (
	KeyLastModified = "wikinote.bluemir.me/last-modified"
)

type Options struct {
}
type Core struct {
	*Options
	store metadata.Store
}

func init() {
	plugins.Register("last-modified", New, &Options{})
}

func New(o interface{}, store metadata.Store) (plugins.Plugin, error) {
	opt, ok := o.(*Options)
	if !ok {
		return nil, errors.Errorf("option not matched")
	}

	return &Core{opt, store}, nil
}

func (c *Core) FileWriteHook(path string, data []byte) ([]byte, error) {
	if err := c.store.Save(path, KeyLastModified, time.Now().UTC().Format(time.RFC3339)); err != nil {
		return data, err
	}
	return data, nil
}
func (c *Core) Footer(path string) ([]byte, error) {
	value, err := c.store.Take(path, KeyLastModified)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return []byte{}, err
	}

	return []byte("last update: " + t.Local().Format(time.RFC3339)), nil
}
