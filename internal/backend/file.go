package backend

import (
	"io"

	"github.com/bluemir/wikinote/internal/backend/files"
	"github.com/bluemir/wikinote/internal/events"
)

func (backend *Backend) FileRead(path string) ([]byte, error) {
	// TODO post hook?
	return backend.files.Read(path)
}
func (backend *Backend) FileReadStream(path string) (io.ReadSeekCloser, error) {
	return backend.files.ReadStream(path)
}

func (backend *Backend) FileWrite(path string, data []byte) error {
	data, err := backend.Plugin.TriggerFileWriteHook(path, data)
	if err != nil {
		backend.hub.Fire(events.Event[Message]{
			Name: "group/admin",
			Detail: Message{
				Text: err.Error(),
			},
		})
		return err
	}

	return backend.files.Write(path, data)
}

func (backend *Backend) FileDelete(path string) error {
	return backend.files.Delete(path)
}
func (backend *Backend) FileList(path string) ([]files.FileInfo, error) {
	return backend.files.List(path)
}
