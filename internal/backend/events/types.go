package events

type FileWritten struct {
	Path string
}

type FileDeleted struct {
	Path string
}

type ErrorOccured struct {
	Error error
}

type SystemMessage struct {
	Message string
}
