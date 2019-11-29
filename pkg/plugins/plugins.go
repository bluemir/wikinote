package plugins

type Plugin interface {
}

type PluginFooter interface {
	Footer(path string) ([]byte, error)
}
type PluginReadHook interface {
	FileReadHook(path string, data []byte) ([]byte, error)
}
type PluginWriteHook interface {
	FileWriteHook(path string, data []byte) ([]byte, error)
}
