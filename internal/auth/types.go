package auth

import "regexp"

type Group struct {
	Name string
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
