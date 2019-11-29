package recent

import (
	"time"

	"github.com/pkg/errors"

	"github.com/bluemir/wikinote/pkg/plugins"
)

const (
	KeyLastModified = "wikinote.bluemir.me/last-modified"
)

type Options struct {
}
type Recents struct {
	opt   *Options
	store *plugins.Store
}

func init() {
	plugins.Register("recent-changes", New, &Options{})
}
func New(o interface{}, store *plugins.Store) (plugins.Plugin, error) {
	opt, ok := o.(*Options)
	if !ok {
		return nil, errors.Errorf("option not matched")
	}

	return &Recents{opt, store}, nil
}

func (r *Recents) FileWriteHook(path string, data []byte) ([]byte, error) {
	// register lastModified
	if err := r.store.Save(&plugins.FileAttr{
		Path:  path,
		Key:   KeyLastModified,
		Value: time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		return data, err
	}

	return data, nil
}
func (r *Recents) Footer(path string) ([]byte, error) {
	attr, err := r.store.Take(&plugins.FileAttr{
		Path: path,
		Key:  KeyLastModified,
	})
	if err != nil {
		if plugins.IsNotFound(err) {
			return []byte(""), nil
		}
		return []byte{}, err
	}
	t, err := time.Parse(time.RFC3339, attr.Value)
	if err != nil {
		return []byte{}, err
	}

	return []byte("last update: " + t.Local().Format(time.RFC3339)), nil
}
