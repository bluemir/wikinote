package plugins

type PluginInit func(opt interface{}) (Plugin, error)

type PluginInitDriver struct {
	Init    PluginInit
	Options interface{}
}

var inits = map[string]PluginInitDriver{}

func Register(name string, init PluginInit, opt interface{}) {
	inits[name] = PluginInitDriver{init, opt}
}
