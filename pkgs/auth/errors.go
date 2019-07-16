package auth

import (
	"fmt"
)

type AuthError struct {
	error
	code ErrorCodeNum
}

func Errorf(code ErrorCodeNum, format string, val ...interface{}) error {
	return Error(code, fmt.Errorf(format, val...))
}
func Error(code ErrorCodeNum, err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*AuthError); ok {
		return e
	}
	return &AuthError{err, code}
}

func ErrorCode(err error) ErrorCodeNum {
	if err == nil {
		return ErrNone
	}
	if e, ok := err.(*AuthError); ok {
		return e.code
	}
	return ErrUnknown
}
