package config

import (
	"github.com/bluemir/wikinote/internal/server"
)

type Config struct {
	Backend   BackendCLIOptions
	Server    server.Config
	LogLevel  int
	LogFormat string
}

func NewConfig() *Config {
	return &Config{
		Backend: BackendCLIOptions{
			AdminUsers: map[string]string{},
		},
	}
}

type BackendCLIOptions struct {
	Wikipath   string
	AdminUsers map[string]string
}
