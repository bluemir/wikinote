package auth

import (
	"github.com/pkg/errors"
)

type User struct {
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"unique"`
	Labels Labels `sql:"type:json"`
}
type Token struct {
	ID        uint `gorm:"primary_key"`
	UserName  string
	HashedKey string `json:"-"`
	RevokeKey string `gorm:"unique"`
}

type Context struct {
	Object  map[string]string
	Subject Subject
	Action  string
}
type Subject struct {
	User  *User
	Token *Token
}

var (
	ErrEmptyHeader   = errors.Errorf("empty header")
	ErrWrongEncoding = errors.Errorf("wrong encoding(not base64)")
	ErrBadToken      = errors.Errorf("bad token. token parse failed")
	ErrNotImplements = errors.Errorf("not implements")

	ErrEmptyAccount = errors.Errorf("empty account")
	ErrUnauthorized = errors.Errorf("unauthorized")
)

type Result int

const (
	Error Result = iota
	Reject
	Accept
	NeedAuthn
)
