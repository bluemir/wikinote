package events

var (
	KindFileWritten = "system.file.written"
)

type FileWritten struct {
	Path string
}
