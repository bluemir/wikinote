package auth

import (
	"encoding/base64"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/bluemir/wikinote/pkgs/utils"
)

func New(store *gorm.DB) (*Manager, error) {
	err := store.AutoMigrate(
		&User{},
		&Token{},
		&UserAttr{},
		&TokenAttr{},
	).Error
	if err != nil {
		return nil, err
	}
	return &Manager{store}, nil
}

type Manager struct {
	store *gorm.DB
}

func (m *Manager) HttpAuth(header string) (*Token, error) {
	if header == "" {
		return nil, Errorf(ErrEmptyHeader, "Empty header")
	}
	arr := strings.SplitN(header, " ", 2)
	switch arr[0] {
	case "Basic", "basic", "Token", "token":
		str, err := base64.StdEncoding.DecodeString(arr[1])
		if err != nil {
			return nil, Error(ErrWrongEncoding, err)
		}

		authStr := strings.SplitN(string(str), ":", 2)
		if len(authStr) != 2 {
			return nil, Errorf(ErrBadToken, "Token invaildated")
		}
		return m.Default(authStr[0], authStr[1])
	case "Bearer", "bearer":
		// TODO
		return nil, Errorf(ErrNotImplement, "Not Implements")
	default:
		return nil, Errorf(ErrNotImplement, "Not Implements")
	}
	return nil, nil
}
func (m *Manager) Default(username, unhashedKey string) (*Token, error) {
	if username == "" && unhashedKey == "" {
		return nil, Errorf(ErrEmptyAccount, "EmptyAccount")
	}
	token := &Token{}
	if err := m.store.Where(&Token{
		HashedKey: hash(unhashedKey, salt(username)),
	}).Take(token).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, Errorf(ErrUnauthorized, "token not found")
		}
		return nil, Error(ErrStore, err)
	}
	return token, nil
}
func (m *Manager) Root(username string) (string, error) {
	_, exist, err := m.GetUser(username)
	if err != nil {
		return "", err
	}
	if !exist {
		if err := m.CreateUser(&User{
			Name: username,
		}); err != nil {
			return "", err
		}
	}

	tokens, err := m.ListToken(username)
	if err != nil {
		return "", err
	}
	for _, token := range tokens {
		if err := m.RevokeToken(token.RevokeKey); err != nil {
			return "", err
		}
	}

	key := utils.RandomString(16)
	if _, err := m.IssueToken(username, key); err != nil {
		return "", err
	}

	if err := m.SetUserAttr(username, "rbac/role-root", "true"); err != nil {
		return "", err
	}

	if err := m.SetUserAttr(username, "user.wikinote.io/email", "root@wikinote.io"); err != nil {
		return "", err
	}

	return key, nil
}
