package config

import (
	"github.com/bluemir/wikinote/internal/server"
)

type Config struct {
	Backend struct {
		Wikipath         string
		VolatileDatabase bool
	}
	Server    server.Config
	LogLevel  int
	LogFormat string
}

func NewConfig() *Config {
	return &Config{}
}
