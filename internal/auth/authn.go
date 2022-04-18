package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (m *Manager) Default(name, unhashedKey string) (*User, error) {
	user := User{}

	if err := m.db.Where(&User{
		Name: name,
	}).Take(&user).Error; err != nil {
		return nil, ErrUnauthorized
	}

	token := &Token{}
	if err := m.db.Where(&Token{
		Username:  name,
		HashedKey: hash(unhashedKey, salt(name)),
	}).Take(token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	return &user, nil
}
func (m *Manager) HTTP(req *http.Request) (*User, error) {
	return m.HTTPHeaderString(req.Header.Get(HeaderAuthorization))
}
func (m *Manager) HTTPHeaderString(header string) (*User, error) {
	if header == "" {
		logrus.Trace("EmptyHeader")
		return nil, ErrEmptyHeader
	}
	method, data := split2(header)
	switch strings.ToLower(method) {
	case "basic", "token":
		str, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			logrus.Error(err)
			return nil, ErrWrongEncoding
		}

		username, key := split2(string(str))
		return m.Default(username, key)
	case "Bearer", "bearer":
		// TODO
		return nil, ErrNotImplements
	default:
		return nil, ErrNotImplements
	}
}
func split2(str string) (string, string) {
	arr := strings.SplitN(str, " ", 2)
	if len(arr) < 2 {
		return arr[0], ""
	}
	return arr[0], arr[1]
}
