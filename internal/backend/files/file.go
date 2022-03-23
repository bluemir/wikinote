package files

import (
	"io"
)

type FileStore interface {
	Read(path string) ([]byte, error)
	ReadStream(path string) (io.ReadSeekCloser, error)
	Write(path string, data []byte) error
	WriteStream(path string, r io.ReadCloser) error
	Delete(path string) error
}
