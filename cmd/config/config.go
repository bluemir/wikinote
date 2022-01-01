package config

import (
	"github.com/bluemir/wikinote/internal/backend"
	"github.com/bluemir/wikinote/internal/server"
)

type Config struct {
	Backend   backend.Config
	Server    server.Config
	LogLevel  int
	LogFormat string
}

func NewConfig() *Config {
	return &Config{
		Backend: backend.InitConfig(),
	}
}
