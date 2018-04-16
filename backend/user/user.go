package user

import (
	"crypto"
	_ "crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name  string
	Email string
	Role  string
}
type Token struct {
	gorm.Model
	UserID    uint
	HashedKey string
}

func hash(username string, key string) string {
	salt := rawHashHex(key + "__salt__" + username)
	return rawHashBase64(salt[:64] + key + salt[64:])
}
func rawHashBase64(input string) string {
	h := crypto.SHA512.New()
	io.WriteString(h, input)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
func rawHashHex(str string) string {
	h := crypto.SHA512.New()
	io.WriteString(h, str)
	return hex.EncodeToString(h.Sum(nil))
}
