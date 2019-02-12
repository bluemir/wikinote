package publish

import (
	"bytes"

	"github.com/bluemir/wikinote/plugins"
	"github.com/pkg/errors"
)

const (
	AttrKey = "publish.plugin.bluemir.me/state"

	stateDraft     = "draft"
	statePublished = "published"
)

func init() {
	plugins.Register("publish", New)
}

func New(opts map[string]string, store plugins.FileAttrStore, auth plugins.AuthManager) plugins.Plugin {
	state, ok := opts["default-state"]
	if ok && state == statePublished {
		return &Publish{
			defaultState: statePublished,
			auth:         auth,
		}
	}
	return &Publish{
		defaultState: stateDraft,
		auth:         auth,
	}
}

type Publish struct {
	defaultState string
	auth         plugins.AuthManager
}

func (publish *Publish) TryRead(path string, user interface{}, attr plugins.FileAttr) error {
	state, err := attr.Get(AttrKey)
	if err != nil {
		state = publish.defaultState
	}
	switch state {
	case stateDraft:
		if publish.auth.Is(user).NotAllow("edit") {
			return errors.Errorf("This is Draft")
		}
	case statePublished:
		return nil
	default:
	}
	return nil
}
func (publish *Publish) TryWrite(path string, user interface{}, attr plugins.FileAttr) error {
	return nil
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
