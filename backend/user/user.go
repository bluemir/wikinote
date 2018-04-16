package user

import (
	"bytes"
	"crypto"
	_ "crypto/sha512"
	"encoding/base64"
	"encoding/hex"

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
	return rawHashBase64(salt[:512] + key + salt[512:])
}
func rawHashBase64(input string) string {
	buf := bytes.NewBuffer([]byte{})
	e := base64.NewEncoder(base64.StdEncoding, buf)
	e.Write(crypto.SHA512.New().Sum([]byte(input)))
	e.Close()

	return buf.String()
}
func rawHashHex(str string) string {
	hashed := crypto.SHA512.New().Sum([]byte(str))
	return hex.EncodeToString(hashed)
}
