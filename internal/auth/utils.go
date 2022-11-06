package auth

import (
	"crypto"
	"encoding/hex"
	"io"
	"strings"
)

func hashRawHex(str string) string {
	hashed := crypto.SHA512.New()
	io.WriteString(hashed, str)
	return hex.EncodeToString(hashed.Sum(nil)) // return 128 length string
}

func hash(str string, salt ...string) string {
	return hashRawHex(strings.Join(append([]string{str}, salt...), "/"))
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
func unique(strs []string) []string {
	set := map[string]struct{}{}
	for _, str := range strs {
		set[str] = struct{}{}
	}
	ret := []string{}
	for str := range set {
		ret = append(ret, str)
	}
	return ret
}
