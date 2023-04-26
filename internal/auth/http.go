package auth

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"github.com/rs/xid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	headerAuthorization   = "Authorization"
	headerWWWAuthenticate = "WWW-Authenticate"
)

func HttpRealm(relam string) string {
	return `Basic realm="` + relam + `"`
}
func LoginHeader(req *http.Request) (string, string) {
	return headerWWWAuthenticate, "basic realm=" + req.URL.Host
}
func (m *Manager) HTTP(req *http.Request) (*User, error) {
	return m.HTTPHeaderString(req.Header.Get(headerAuthorization))
}
func (m *Manager) HTTPHeaderString(header string) (*User, error) {
	if header == "" {
		logrus.Trace("EmptyHeader")
		return nil, ErrUnauthorized
	}
	method, data := split2(header, " ")
	switch strings.ToLower(method) {
	case "basic", "token":
		str, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			logrus.Error(err)
			return nil, ErrUnauthorized // maybe bad request?
		}

		username, key := split2(string(str), ":")
		return m.Default(username, key)
	case "Bearer", "bearer":
		return nil, ErrNotImplements
	default:
		return nil, ErrNotImplements
	}
}
func split2(str string, sep string) (string, string) {
	arr := strings.SplitN(str, sep, 2)
	if len(arr) < 2 {
		return arr[0], ""
	}
	return arr[0], arr[1]
}

func (m *Manager) NewHTTPToken(username string, expireAt time.Time) (string, error) {
	user := &User{}
	if err := m.db.Where(&User{Name: username}).Take(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", ErrUnauthorized
		}
		return "", err
	}

	newKey := hash(xid.New().String(), user.Salt)

	if _, err := m.IssueToken(username, newKey, &expireAt); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(strings.Join([]string{user.Name, newKey}, ":"))), nil
}

func (m *Manager) RevokeHTTPToken(req *http.Request) error {
	username, key, err := parseHTTPHeader(req.Header.Get(headerAuthorization))
	if err != nil {

		return err
	}
	return m.RevokeToken(username, key)
}

func parseHTTPHeader(header string) (string, string, error) {
	method, data := split2(header, " ")

	switch strings.ToLower(method) {
	case "basic", "token", "bearer":
		c, err := base64.StdEncoding.DecodeString(strings.TrimSpace(data))
		if err != nil {
			return "", "", err
		}

		name, key := split2(string(c), ":")

		return name, key, nil
	default:
		return "", "", ErrUnauthorized // unknown method
	}
}
