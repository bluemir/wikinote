package backend

import (
	"fmt"

	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	"github.com/bluemir/wikinote/pkgs/auth"
	"github.com/bluemir/wikinote/pkgs/config"
	"github.com/bluemir/wikinote/pkgs/fileattr"
	"github.com/bluemir/wikinote/plugins"
)

type pluginList struct {
	footer         []plugins.FooterPlugin
	postSave       []plugins.PostSavePlugin
	preSave        []plugins.PreSavePlugin
	onReadWiki     []plugins.ReadWikiPlugin
	authz          []plugins.AuthzPlugin
	registerRouter map[string]plugins.RegisterRouterPlugin
}

func loadPlugins(conf *config.Config, store fileattr.Store, authManager *auth.Manager) (*pluginList, error) {
	// TODO can on/off
	pl := &pluginList{
		footer:         []plugins.FooterPlugin{},
		postSave:       []plugins.PostSavePlugin{},
		preSave:        []plugins.PreSavePlugin{},
		onReadWiki:     []plugins.ReadWikiPlugin{},
		authz:          []plugins.AuthzPlugin{},
		registerRouter: map[string]plugins.RegisterRouterPlugin{},
	}

	for name, pconf := range conf.Plugins {
		pc, ok := pconf.(map[string]interface{})
		if !ok {
			pc = map[string]interface{}{}
		}

		p, err := plugins.New(name, flat(pc), store, authManager) // TODO config
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		logrus.Debugf("plugin '%s' is initialize", name)
		if f, ok := p.(plugins.FooterPlugin); ok {
			logrus.Debugf("footer plugin '%s'", name)
			pl.footer = append(pl.footer, f)
		}
		if plugin, ok := p.(plugins.PreSavePlugin); ok {
			logrus.Debugf("pre save plugin '%s'", name)
			pl.preSave = append(pl.preSave, plugin)
		}
		if a, ok := p.(plugins.PostSavePlugin); ok {
			logrus.Debugf("post save plugin '%s'", name)
			pl.postSave = append(pl.postSave, a)
		}
		if plugin, ok := p.(plugins.ReadWikiPlugin); ok {
			logrus.Debugf("read plugin '%s'", name)
			pl.onReadWiki = append(pl.onReadWiki, plugin)
		}
		if plugin, ok := p.(plugins.AuthzPlugin); ok {
			logrus.Debugf("permission plugin '%s'", name)
			pl.authz = append(pl.authz, plugin)
		}
		if a, ok := p.(plugins.RegisterRouterPlugin); ok {
			logrus.Debugf("resiger route plugin '%s'", name)
			pl.registerRouter[name] = a
		}
	}

	for _, pconf := range conf.PluginsV2 {
		buf, err := yaml.Marshal(pconf.Options)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		p, err := plugins.NewV2(pconf.Name, nil, buf)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		name := pconf.Name
		logrus.Debugf("plugin '%s' is initialize", name)
		if f, ok := p.(plugins.FooterPlugin); ok {
			logrus.Debugf("footer plugin '%s'", name)
			pl.footer = append(pl.footer, f)
		}
		if plugin, ok := p.(plugins.PreSavePlugin); ok {
			logrus.Debugf("pre save plugin '%s'", name)
			pl.preSave = append(pl.preSave, plugin)
		}
		if a, ok := p.(plugins.PostSavePlugin); ok {
			logrus.Debugf("post save plugin '%s'", name)
			pl.postSave = append(pl.postSave, a)
		}
		if plugin, ok := p.(plugins.ReadWikiPlugin); ok {
			logrus.Debugf("read plugin '%s'", name)
			pl.onReadWiki = append(pl.onReadWiki, plugin)
		}
		if plugin, ok := p.(plugins.AuthzPlugin); ok {
			logrus.Debugf("permission plugin '%s'", name)
			pl.authz = append(pl.authz, plugin)
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
