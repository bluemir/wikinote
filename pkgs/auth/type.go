package auth

type User struct {
	ID   string `gorm:"primary_key"`
	Name string `gorm:"unique"`
}
type Token struct {
	ID        string `gorm:"primary_key"`
	UserName  string
	HashedKey string `json:"-"`
	RevokeKey string `gorm:"unique"`
}

// AuthAttr
type Attr struct {
	ID        string `gorm:"primary_key"`
	Namespace string `gorm:"primary_key"`
	Key       string `gorm:"primary_key"`
	Value     string
}

// helper
type Subject interface {
	Attr(key string) string
}
type Object interface {
	Attr(key string) string
}
type Context struct {
	Subject Subject
	Object  Object
	Action  string
}
type Result int

const (
	Unknown Result = iota
	Reject
	Accept
)

// ABAC(Context{subject, object, action}) == Reject
