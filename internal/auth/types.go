package auth

type User struct {
	ID     uint   `gorm:"primary_key"`
	Name   string `gorm:"unique"`
	Labels Labels `sql:"type:json"`
	Salt   string
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
type RoleBinding struct {
	Username string
	Rolename string
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
func (kvs KeyValues) isSubsetOf(resource Resource) bool {
	for k, v := range kvs {
		if resource.Get(k) != v {
			return false
		}
	}
	return true
}
