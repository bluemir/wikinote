package context

import (
	"github.com/docker/libkv/store"

	"github.com/bluemir/wikinote/backend/config"
)

type Core interface {
	Store() store.Store
	Config() config.Config
}
