package plugins

import (
	"html/template"

	"github.com/bluemir/wikinote/internal/backend/metadata"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

type PluginConfig struct {
	Name    string
	Options interface{}
}

func New(configs []PluginConfig, store metadata.Store) (*Manager, error) {
	manager := &Manager{Route: map[string]PluginRoute{}}
	for _, conf := range configs {
		logrus.Infof("Initialize plugin: %s", conf.Name)
		p, ok := inits[conf.Name]
		if !ok {
			return nil, errors.Errorf("plugin not found")
		}

		if conf.Options != nil {
			buf, err := yaml.Marshal(conf.Options)
			if err != nil {
				return nil, err
			}

			logrus.Infof("%s", buf)
			if err := yaml.Unmarshal(buf, p.Options); err != nil {
				return nil, err
			}
		}

		logrus.Infof("%#v", p.Options)
		plugin, err := p.Init(p.Options, store)
		if err != nil {
			return nil, err
		}

		if v, ok := plugin.(PluginFooter); ok {
			manager.Footer = append(manager.Footer, v)
			logrus.Tracef("footer detected")
		}
		if v, ok := plugin.(PluginWriteHook); ok {
			manager.WriteHook = append(manager.WriteHook, v)
			logrus.Tracef("writeHook detected")
		}
		if v, ok := plugin.(PluginRoute); ok {
			manager.Route[conf.Name] = v
			logrus.Tracef("route detected")
		}
	}
	return manager, nil
}

type Manager struct {
	Footer    []PluginFooter
	WriteHook []PluginWriteHook
	Route     map[string]PluginRoute
}

func (m *Manager) TriggerFileReadHook(path string, data []byte) ([]byte, error) {
	return data, nil
}
func (m *Manager) TriggerFileWriteHook(path string, data []byte) ([]byte, error) {
	for _, hook := range m.WriteHook {
		newData, err := hook.FileWriteHook(path, data)
		if err != nil {
			// just log it and skip this hook
			// TODO report to event store and show it to admin?
			logrus.Error(err)
		}
		data = newData
	}
	return data, nil
}
func (m *Manager) GetWikiFooter(path string) ([]template.HTML, error) {
	result := []template.HTML{}
	for _, hook := range m.Footer {
		data, err := hook.Footer(path)
		if err != nil {
			logrus.Error(err)
			result = append(result, template.HTML("error in plugin: "+err.Error()))
			//return result, err
		}
		result = append(result, template.HTML(data))
	}
	return result, nil
}
func (m *Manager) RouteHook(app gin.IRouter) error {
	for name, route := range m.Route {
		if err := route.Route(app.Group(name)); err != nil {
			return err
		}
	}
	return nil
}
