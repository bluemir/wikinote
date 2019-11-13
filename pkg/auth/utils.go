package auth

import (
	"crypto"
	"encoding/hex"
	"io"
)

func hashRawHex(str string) string {
	hashed := crypto.SHA512.New()
	io.WriteString(hashed, str)
	return hex.EncodeToString(hashed.Sum(nil))
}
func hash(unhashedKey, saltSeed string) string {
	salt := hashRawHex(saltSeed)
	return hashRawHex(salt[:64] + unhashedKey + salt[64:])
}
func salt(username string) string {
	return "__salt__" + username + "__salt__"
}

const (
	HeaderAuthorization   = "Authorization"
	HeaderWWWAuthenticate = "WWW-Authenticate"
)

func HttpRealm(relam string) string {
	return `Basic realm="` + relam + `"`
}
