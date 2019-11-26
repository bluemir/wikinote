package plugins

import yaml "gopkg.in/yaml.v2"

type PluginInit func(opt interface{}) Plugin

type Plugin interface {
}

type PluginFooter interface {
	Footer(path string) error
}
type PluginReadHook interface {
	FileReadHook(path string, data []byte) ([]byte, error)
}
type PluginWriteHook interface {
	FileWriteHook(path string, data []byte) ([]byte, error)
}

type PluginInitWithOption struct {
	Init    PluginInit
	Options interface{}
}

var inits = map[string]PluginInitWithOption{}

func Register(name string, init PluginInit, opt interface{}) {
	inits[name] = PluginInitWithOption{init, opt}
}

func NewManager(configs []PluginConfig) (*Manager, error) {
	manager := &Manager{}
	for _, conf := range configs {
		buf, err := yaml.Marshal(conf.Options)
		if err != nil {
			return nil, err
		}
		p := inits[conf.Name]

		if err := yaml.Unmarshal(buf, p.Options); err != nil {
			return nil, err
		}

		plugin := p.Init(p.Options)

		if v, ok := plugin.(PluginFooter); ok {
			manager.Footer = append(manager.Footer, v)
		}
	}
	return manager, nil
}

type Manager struct {
	Footer []PluginFooter
}

func (m *Manager) TriggerFileReadHook(path string, data []byte) ([]byte, error) {
	return nil, nil
}

type PluginConfig struct {
	Name    string
	Options interface{}
}
