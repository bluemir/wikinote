package backend

import (
	"path/filepath"
	"strings"

	"github.com/bluemir/wikinote/internal/fileattr"
)

func (backend *Backend) Object(path string) (map[string]string, error) {
	if strings.HasPrefix(path, "/!/") {
		return map[string]string{"kind": "special"}, nil
	}
	// TODO make kind from ext
	attrs, err := backend.FileAttr.Find(&fileattr.FileAttr{
		Path: path,
	})

	if err != nil {
		return map[string]string{}, err
	}

	result := map[string]string{}
	for _, attr := range attrs {
		result[attr.Key] = attr.Value
	}

	ext := filepath.Ext(path)
	switch ext {
	case ".md":
		result["kind"] = "wiki"
	}

	return result, nil
}

func (backend *Backend) FileRead(path string) ([]byte, error) {
	return backend.files.Read(path)
}

func (backend *Backend) FileWrite(path string, data []byte) error {
	data, err := backend.Plugin.TriggerFileWriteHook(path, data)
	if err != nil {
		return err
	}

	return backend.files.Write(path, data)
}

func (backend *Backend) FileDelete(path string) error {
	return backend.files.Delete(path)
}
