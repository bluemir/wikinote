package auth

import (
	"regexp"
	"time"
)

type User struct {
	ID     uint   `gorm:"primary_key" json:"-"`
	Name   string `gorm:"unique" json:"name"`
	Groups Set    `gorm:"type:bytes;serializer:gob" json:"groups"`
	Labels Labels `gorm:"type:bytes;serializer:gob" json:"labels" expr:"labels"`
	Salt   string `json:"-"`
}
type Group struct {
	Name string `json:"name,omitempty"`
}
type Token struct {
	Id        uint   `gorm:"primary_key" json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	HashedKey string `json:"-,omitempty"`
	RevokeKey string `json:"revoke_key,omitempty"`
	ExpiredAt *time.Time
}
type Role struct {
	Rules []Rule `yaml:"rules"`
}
type Rule struct {
	Verbs     []Verb          `yaml:"verbs"`
	Resources []ResourceMatch `yaml:"resources"`
}
type Resource interface {
	Get(key string) string
	KeyValues() KeyValues
}
type Verb string

type KeyValues map[string]string

func (kvs KeyValues) Get(key string) string {
	return kvs[key]
}
func (kvs KeyValues) KeyValues() KeyValues {
	return kvs
}

type ResourceMatch map[string]*Regexp

type Regexp struct {
	*regexp.Regexp
}

// UnmarshalText unmarshals json into a regexp.Regexp
func (r *Regexp) UnmarshalText(b []byte) error {
	regex, err := regexp.Compile(string(b))
	if err != nil {
		return err
	}

	r.Regexp = regex

	return nil
}

// MarshalText marshals regexp.Regexp as string
func (r *Regexp) MarshalText() ([]byte, error) {
	if r.Regexp != nil {
		return []byte(r.Regexp.String()), nil
	}

	return nil, nil
}

type Subject struct {
	Kind Kind   `gorm:"primaryKey;size:256" expr:"kind"`
	Name string `gorm:"primaryKey;size:256" expr:"name"`
}

type Kind string

const (
	KindUser  Kind = "user"
	KindGroup Kind = "group"
)

type TokenOpt func(*Token)

func ExpiredAt(t time.Time) func(*Token) {
	return func(token *Token) {
		token.ExpiredAt = &t
	}
}
func ExpiredAfter(d time.Duration) func(*Token) {
	return func(token *Token) {
		t := time.Now().Add(d)
		token.ExpiredAt = &t
	}
}

type Labels map[string]string
type List []string
type Set map[string]struct{}
