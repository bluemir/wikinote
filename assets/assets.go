package assets

import (
	"io/fs"
	"path"
	"strings"
)

var (
	HTMLTemplates fs.FS
	Static        fs.FS
)

const (
	src  = "src"
	dist = "dist"
)

// mappedFS maps js/* -> {prefix}/js/*, css/* -> {prefix}/css/*
type mappedFS struct {
	fs     fs.FS
	prefix string // "dist" for prod, "src" for dev
}

func (m *mappedFS) Open(name string) (fs.File, error) {
	if strings.HasPrefix(name, "js/") || strings.HasPrefix(name, "css/") {
		return m.fs.Open(path.Join(m.prefix, name))
	}
	return m.fs.Open(name) // lib/*, images/*
}
