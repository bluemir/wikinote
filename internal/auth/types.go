package auth

import "strings"

type User struct {
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"unique"`
	Labels Labels `sql:"type:json"`
	Salt   string
}

func (user User) Roles() []string {
	result := []string{}
	for k := range user.Labels {
		if strings.HasPrefix(k, "role/") {
			result = append(result, strings.TrimPrefix(k, "role/"))
		}
	}
	return result
}
func (user *User) AddRole(role string) {
	user.Labels["role/"+role] = "true"
}
func (user *User) RemoveRole(role string) {
	delete(user.Labels, "role/"+role)
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
	Resource ResourceExprs `yaml:"resource"`
	Verbs    []Verb        `yaml:"verbs"`
	Expr     RuleExpr      `yaml:"expr"`
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
