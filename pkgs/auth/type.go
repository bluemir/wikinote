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

type Attr struct {
	ID    string
	Key   string
	Value string
}

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
