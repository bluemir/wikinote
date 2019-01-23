package backend

import (
	"fmt"

	"github.com/bluemir/wikinote/pkgs/config"
	"github.com/bluemir/wikinote/plugins"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
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
	afterWikiSave  []plugins.AfterWikiSavePlugin
	registerRouter map[string]plugins.RegisterRouterPlugin
}

func loadPlugins(db *gorm.DB, conf *config.Config) (*pluginList, error) {
	// TODO can on/off
	pl := &pluginList{
		footer:         []plugins.FooterPlugin{},
		afterWikiSave:  []plugins.AfterWikiSavePlugin{},
		registerRouter: map[string]plugins.RegisterRouterPlugin{},
	}

	for name, pconf := range conf.Plugins {
		pc, ok := pconf.(map[string]interface{})
		if !ok {
			pc = map[string]interface{}{}
		}

		p, err := plugins.New(name, db, flat(pc)) // TODO config
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		logrus.Debugf("plugin '%s' is initialize", name)
		if f, ok := p.(plugins.FooterPlugin); ok {
			pl.footer = append(pl.footer, f)
		}
		if a, ok := p.(plugins.AfterWikiSavePlugin); ok {
			pl.afterWikiSave = append(pl.afterWikiSave, a)
		}
		if a, ok := p.(plugins.RegisterRouterPlugin); ok {
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
