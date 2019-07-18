package backend

import (
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"

	"github.com/bluemir/wikinote/pkgs/config"
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
type pluginsCore struct {
	*backend
}

func (core *pluginsCore) Auth() plugins.CoreAuth {
	return core.backend.Auth()
}
func (core *pluginsCore) File() plugins.CoreFile {
	return core
}
func (core *pluginsCore) Attr() plugins.CoreFileAttr {
	return core.backend.File().AttrStore()
}

func (b *backend) loadPlugins(conf *config.Config) error {
	pl := &pluginList{
		footer:         []plugins.FooterPlugin{},
		postSave:       []plugins.PostSavePlugin{},
		preSave:        []plugins.PreSavePlugin{},
		onReadWiki:     []plugins.ReadWikiPlugin{},
		authz:          []plugins.AuthzPlugin{},
		registerRouter: map[string]plugins.RegisterRouterPlugin{},
	}
	core := &pluginsCore{b}

	for _, pconf := range conf.Plugins {
		buf, err := yaml.Marshal(pconf.Options)
		if err != nil {
			logrus.Error(err)
			return err
		}
		p, err := plugins.New(pconf.Name, core, buf)
		if err != nil {
			logrus.Error(err)
			return err
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
	b.plugins = pl
	return nil
}
