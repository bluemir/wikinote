package backend

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/pkgs/config"
	"github.com/bluemir/wikinote/pkgs/fileattr"
	"github.com/bluemir/wikinote/plugins"
)

const (
	defaultRule = `
rules:
  admin:  [ "view", "edit", "user", "search" ]
  editor: [ "view', "edit", "attach", "search" ]
  viewer: [ "view", "search" ]
  guest:  [ "view" ]
`
)

type pluginList struct {
	footer         []plugins.FooterPlugin
	postSave       []plugins.PostSavePlugin
	registerRouter map[string]plugins.RegisterRouterPlugin
}

func loadPlugins(conf *config.Config, store fileattr.Store) (*pluginList, error) {
	// TODO can on/off
	pl := &pluginList{
		footer:         []plugins.FooterPlugin{},
		postSave:       []plugins.PostSavePlugin{},
		registerRouter: map[string]plugins.RegisterRouterPlugin{},
	}

	for name, pconf := range conf.Plugins {
		pc, ok := pconf.(map[string]interface{})
		if !ok {
			pc = map[string]interface{}{}
		}

		p, err := plugins.New(name, flat(pc), store) // TODO config
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		logrus.Debugf("plugin '%s' is initialize", name)
		if f, ok := p.(plugins.FooterPlugin); ok {
			logrus.Debugf("footer plugin '%s'", name)
			pl.footer = append(pl.footer, f)
		}
		if a, ok := p.(plugins.PostSavePlugin); ok {
			logrus.Debugf("post save plugin '%s'", name)
			pl.postSave = append(pl.postSave, a)
		}
		if a, ok := p.(plugins.RegisterRouterPlugin); ok {
			logrus.Debugf("resiger route plugin '%s'", name)
			pl.registerRouter[name] = a
		}
	}
	return pl, nil
}
func flat(conf map[string]interface{}) map[string]string {
	result := map[string]string{}
	for k, v := range conf {
		switch val := v.(type) {
		case map[string]interface{}:
			res := flat(val)
			for resk, resv := range res {
				result[k+"."+resk] = resv
			}
		default:
			result[k] = fmt.Sprint(v)
		}
	}
	return result
}
