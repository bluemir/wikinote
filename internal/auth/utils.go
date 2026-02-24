package auth

import (
	"crypto"
	"encoding/hex"
	"io"
	"strings"

	"golang.org/x/crypto/bcrypt"
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

func bcryptHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func bcryptCompare(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
