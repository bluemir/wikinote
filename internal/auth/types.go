package auth

import (
	"github.com/pkg/errors"
)

var (
	ErrEmptyHeader   = errors.Errorf("empty header")
	ErrWrongEncoding = errors.Errorf("wrong encoding(not base64)")
	ErrBadToken      = errors.Errorf("bad token. token parse failed")
	ErrNotImplements = errors.Errorf("not implements")

	ErrEmptyAccount = errors.Errorf("empty account")
	ErrUnauthorized = errors.Errorf("unauthorized")
)

type User struct {
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"unique"`
	Labels Labels `sql:"type:json"`
}
type Token struct {
	ID        uint `gorm:"primary_key"`
	Username  string
	HashedKey string `json:""`
	RevokeKey string
}
type Role struct {
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules"`
}
type Rule struct {
	Resource KeyValues `yaml:"resource"`
	Verbs    []Verb    `yaml:"verb"`
}
type RoleBinding struct {
	Username string
	Rolename string
}
type Resource interface {
	Get(key string) string
}
type Verb string

type KeyValues map[string]string

func (kvs KeyValues) Get(key string) string {
	return kvs[key]
}
func (kvs KeyValues) isSubsetOf(resource Resource) bool {
	for k, v := range kvs {
		if resource.Get(k) != v {
			return false
		}
	}
	return true
}
