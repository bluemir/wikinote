package buildinfo

import (
	"crypto"
	"encoding/hex"
	"io"
)

var (
	Version   string
	AppName   string
	BuildTime string
)

func Signature() string {
	hashed := crypto.SHA512.New()

	io.WriteString(hashed, AppName)
	io.WriteString(hashed, Version)
	io.WriteString(hashed, BuildTime)

	return hex.EncodeToString(hashed.Sum(nil))
}
