package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

func (m *Manager) Default(name, unhashedKey string) (*Token, error) {
	if name == "" && unhashedKey == "" {
		return nil, ErrEmptyAccount
	}

	token := &Token{}
	if err := m.db.Where(&Token{
		UserName:  name,
		HashedKey: hash(unhashedKey, salt(name)),
	}).Take(token).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	return token, nil
}
func (m *Manager) HTTP(header http.Header) (*Token, error) {
	return m.HTTPHeaderString(header.Get(HeaderAuthorization))
}
func (m *Manager) HTTPHeaderString(header string) (*Token, error) {
	if header == "" {
		logrus.Trace("EmptyHeader")
		return nil, ErrEmptyHeader
	}
	arr := strings.SplitN(header, " ", 2)
	switch arr[0] {
	case "Basic", "basic", "Token", "token":
		str, err := base64.StdEncoding.DecodeString(arr[1])
		if err != nil {
			logrus.Error(err)
			return nil, ErrWrongEncoding
		}

		authStr := strings.SplitN(string(str), ":", 2)
		if len(authStr) != 2 {
			return nil, ErrBadToken
		}
		return m.Default(authStr[0], authStr[1])
	case "Bearer", "bearer":
		// TODO
		return nil, ErrNotImplements
	default:
		return nil, ErrNotImplements
	}
}
