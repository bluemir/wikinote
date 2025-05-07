package backend

import (
	"context"
	"io"
	"io/fs"

	"github.com/bluemir/wikinote/internal/backend/events"
	"github.com/bluemir/wikinote/internal/backend/files"
)

func (backend *Backend) FileRead(path string) ([]byte, error) {
	// TODO post hook?
	return backend.Files.Read(path)
}
func (backend *Backend) FileReadStream(path string) (io.ReadSeekCloser, fs.FileInfo, error) {
	return backend.Files.ReadStream(path)
}

func (backend *Backend) FileWrite(path string, data []byte) error {
	defer backend.hub.Publish(context.Background(), events.FileWritten{
		Path: path,
	})

	return backend.Files.Write(path, data)
}
func (backend *Backend) FileWriteStream(path string, reader io.Reader) error {
	defer backend.hub.Publish(context.Background(), events.FileWritten{
		Path: path,
	})

	return backend.Files.WriteStream(path, reader)
}

func (backend *Backend) FileDelete(path string) error {
	defer backend.hub.Publish(context.Background(), events.FileDeleted{Path: path})
	return backend.Files.Delete(path)
}
func (backend *Backend) FileList(path string) ([]files.FileInfo, error) {
	return backend.Files.List(path)
}
func (backend *Backend) FileMove(oldPath, newPath string) error {
	return backend.Files.Move(oldPath, newPath)
}
