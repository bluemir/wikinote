package files

type FileStore interface {
	Read(path string) ([]byte, error)
	Write(path string, data []byte) error
	Delete(path string) error
}
