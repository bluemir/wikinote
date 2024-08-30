package auth

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/bluemir/wikinote/internal/datastruct"
	"golang.org/x/exp/maps"
)

type Resource interface {
	Get(key string) string
	KeyValues() KeyValues
}
type Verb string

type KeyValues = datastruct.KeyValues

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

type Labels map[string]string
type List []string
type Set map[string]struct{}

func setFromArray(arr []string) Set {
	s := Set{}
	for _, v := range arr {
		s[v] = struct{}{}
	}
	return s
}

var x = struct{}{}

func (s Set) Add(vs ...string) {
	for _, v := range vs {
		s[v] = x
	}
}
func (s Set) MarshalJSON() ([]byte, error) {
	return json.Marshal(maps.Keys(s))
}
func (s Set) ToArray() []string {
	result := []string{}
	for k := range s {
		result = append(result, k)
	}
	return result
}
func (s Set) String() string {
	return strings.Join(s.ToArray(), ",")
}
