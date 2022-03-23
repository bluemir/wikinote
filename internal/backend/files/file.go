package files

//import "io"

type FileStore interface {
	Read(path string) ([]byte, error)
	//ReadStream(path string) (io.ReadCloser, error)
	Write(path string, data []byte) error
	//WriteStream(io.ReadCloser) error
	Delete(path string) error
}
