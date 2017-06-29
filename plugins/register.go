package plugins

type Plugin interface {
}

var plugins = map[string]PluginInit{}

type PluginInit func(opts map[string]string) Plugin

func Register(name string, initFunc PluginInit) {
	plugins[name] = initFunc
}

func New(name string, opts map[string]string) (Plugin, error) {
	return plugins[name](opts), nil
}
func List() []string {
	//return names of plugin
	result := []string{}
	for name, _ := range plugins {
		result = append(result, name)
	}
	return result
}
