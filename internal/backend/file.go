package backend

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/bluemir/wikinote/internal/backend/files"
)

func (backend *Backend) FileRead(path string) ([]byte, error) {
	// TODO post hook?
	return backend.files.Read(path)
}
func (backend *Backend) FileReadStream(path string) (io.ReadSeekCloser, fs.FileInfo, error) {
	return backend.files.ReadStream(path)
}

func (backend *Backend) FileWrite(path string, data []byte) error {
	data, err := backend.Plugin.TriggerFileWriteHook(path, data)
	if err != nil {
		backend.hub.Publish("system", Message{
			Text: err.Error(),
		})
		return err
	}

	defer backend.hub.Publish("system.file.written", Message{
		Text: fmt.Sprintf("file written,path=%s", path),
	})

	return backend.files.Write(path, data)
}

func (backend *Backend) FileDelete(path string) error {
	defer backend.hub.Publish("system.file.deleted", Message{
		Text: fmt.Sprintf("file deleted,path=%s", path),
	})
	return backend.files.Delete(path)
}
func (backend *Backend) FileList(path string) ([]files.FileInfo, error) {
	return backend.files.List(path)
}
func (backend *Backend) FileMove(oldPath, newPath string) error {
	return backend.files.Move(oldPath, newPath)
}
