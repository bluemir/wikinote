package auth

type ErrorCodeNum int32

const (
	ErrNone ErrorCodeNum = iota
	ErrEmptyAccount
	ErrWrongEncoding
	ErrEmptyHeader
	ErrBadToken
	ErrNotImplement
	ErrStore
	ErrUnauthorized
	ErrUnknown
)
