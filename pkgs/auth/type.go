package auth

type User struct {
	ID   uint   `gorm:"primary_key"`
	Name string `gorm:"unique"`
}
type Token struct {
	ID        uint `gorm:"primary_key"`
	UserName  string
	HashedKey string `json:"-"`
	RevokeKey string `gorm:"unique"`
}

// AuthAttr
type Attr struct {
	Key   string
	Value string
}
type UserAttr struct {
	UserId uint   `gorm:"primary_key;auto_increment:false"`
	Key    string `gorm:"primary_key;auto_increment:false"`
	Value  string
}
type TokenAttr struct {
	TokenId uint   `gorm:"primary_key;auto_increment:false"`
	Key     string `gorm:"primary_key;auto_increment:false"`
	Value   string
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
