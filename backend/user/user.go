package user

import (
	"bytes"
	"crypto"
	_ "crypto/sha512"
	"encoding/base64"
)

type User struct {
	Id       string
	Password *Password
	Email    string
	Role     string
}

type Password struct {
	Data string
	Salt string
}

func NewPassword(origin string) *Password {
	p := &Password{}
	p.Set(origin)
	return p
}
func (p *Password) Set(b string) {
	p.Salt = RandomString(12)
	p.Data = encodePassword(b, p.Salt)
}

func (p *Password) Check(b string) bool {
	if p == nil {
		return false
	}
	return p.Data == encodePassword(b, p.Salt)
}

func encodePassword(origin, salt string) string {
	buf := bytes.NewBuffer([]byte{})
	e := base64.NewEncoder(base64.StdEncoding, buf)
	e.Write(crypto.SHA512.New().Sum([]byte(origin + salt)))
	e.Close()

	return buf.String()
}
