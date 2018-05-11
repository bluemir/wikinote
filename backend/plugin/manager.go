package plugin

import (
	"fmt"
	"html/template"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/plugins"
)

type Manager interface {
	Footer(path string) []PluginResult
}
type PluginResult struct {
	Data template.HTML
	Err  error
}

func New(db *gorm.DB, pluginsConf map[interface{}]interface{}) (Manager, error) {
	pluginNames := plugins.List()

	pm := &pluginManager{}

	for _, name := range pluginNames {
		logrus.Debugf("pluginconf %+v, %+v %s", pluginsConf, pluginsConf[name], name)
		pc, ok := pluginsConf[name].(map[interface{}]interface{})
		if !ok {
			pc = map[interface{}]interface{}{}
		}

		p, err := plugins.New(name, db, flat(pc)) // TODO config
		//p, err := plugins.New(name, db, flat(pc)) // TODO config
		if err != nil {
			logrus.Error(err)
		}
		logrus.Debugf("plugin '%s' is initialize", name)
		if f, ok := p.(plugins.FooterPlugin); ok {
			pm.footer = append(pm.footer, f)
		}
	}

	return pm, nil
}
func flat(conf map[interface{}]interface{}) map[string]string {
	result := map[string]string{}
	for k, v := range conf {
		switch val := v.(type) {
		case map[interface{}]interface{}:
			res := flat(val)
			for resk, resv := range res {
				result[k.(string)+"."+resk] = resv
			}
		default:
			result[k.(string)] = fmt.Sprint(v)
		}
	}
	return result
}

type pluginManager struct {
	footer        []plugins.FooterPlugin
	afterWikiSave []plugins.AfterWikiSavePlugin
}

func (pm *pluginManager) Footer(path string) []PluginResult {
	result := []PluginResult{}
	for _, p := range pm.footer {
		d, e := p.Footer(path)
		result = append(result, PluginResult{
			d, e,
		})
	}
	return result
}
func (pm *pluginManager) AfterWikiSave(path string, data []byte) error {
	for _, p := range pm.afterWikiSave {
		e := p.AfterWikiSave(path, data)
		if e != nil {
			return e
		}
	}
	return nil
}
