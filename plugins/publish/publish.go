package publish

import (
	"bytes"

	yaml "gopkg.in/yaml.v3"

	"github.com/bluemir/wikinote/plugins"
)

const (
	AttrKey = "publish.plugin.bluemir.me/state"

	stateDraft     = "draft"
	statePublished = "published"
)

/*
TODO
save to draft. And copy to md when publish
*/

func init() {
	plugins.Register("publish", New)
}

func New(core plugins.Core, confBuf []byte) (plugins.Plugin, error) {
	opts := &Options{
		defaultState: statePublished,
	}

	if err := yaml.Unmarshal(confBuf, opts); err != nil {
		return nil, err
	}
	return &Publish{core, opts}, nil
}

type Publish struct {
	plugins.Core
	opts *Options
}
type Options struct {
	defaultState string
}

func (publish *Publish) OnReadWiki(ctx *plugins.AuthContext, path string, data []byte) ([]byte, error) {
	if ctx.Object.Attr(AttrKey) == stateDraft && ctx.Subject.Attr("publish.plugins.bluemir.me/auth") == "true" {
		return data, nil
	}
	return []byte("it is draft"), nil
}
func (publish *Publish) OnPreSave(path string, data []byte, attr plugins.FileAttr) ([]byte, error) {
	// TODO make to api
	// try to parse first line
	// if matched, set file attribute
	if bytes.HasPrefix(data, []byte(".publish")) {
		err := attr.Set(AttrKey, statePublished)
		if err != nil {
			return nil, err
		}
		return data[len(".publish"):], nil
	}
	if bytes.HasPrefix(data, []byte(".draft")) {
		err := attr.Set(AttrKey, stateDraft)
		if err != nil {
			return nil, err
		}
		return data[len(".draft"):], nil
	}
	return data, nil
}
