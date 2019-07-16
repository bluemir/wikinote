package auth

import (
	"encoding/base64"
	"strings"

	"github.com/bluemir/go-utils/auth/utils"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

func New(store *gorm.DB) (*Manager, error) {
	store.AutoMigrate(
		&User{},
		&Token{},
		&Attr{},
	)
	return &Manager{store}, nil
}

type Manager struct {
	store *gorm.DB
}

func (m *Manager) CreateUser(u *User) error {
	u.ID = xid.New().String()
	if err := m.store.Create(u).Error; err != nil {
		return errors.Wrapf(err, "User already exist")
	}
	return nil
}
func (m *Manager) GetUser(username string) (*User, bool, error) {
	user := &User{}
	if err := m.store.Where("name = ?", username).First(user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return user, true, nil
}
func (m *Manager) ListUser(filter ...string) ([]User, error) {
	users := []User{}
	if err := m.store.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
func (m *Manager) DeleteUser(username string) error {
	if err := m.store.Where("name = ?", username).Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}
func (m *Manager) IssueToken(username, unhashedKey string) (*Token, error) {
	_, ok, err := m.GetUser(username)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.Errorf("user not found")
	}

	id := xid.New().String()
	token := &Token{
		ID:        id,
		UserName:  username,
		HashedKey: hash(unhashedKey, salt(username)),
		RevokeKey: hash(username+id+unhashedKey[:4], "__revoke__"),
	}
	if err := m.store.Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}
func (m *Manager) ListToken(username string) ([]Token, error) {
	result := []Token{}
	err := m.store.Where("username = ?", username).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, err
}
func (m *Manager) RevokeToken(revokeKey string) error {
	err := m.store.Where("revoke_key = ?", revokeKey).Delete(&Token{}).Error
	if err != nil {
		return err
	}
	return nil
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
	user := &User{
		Name: username,
	}
	if err := m.CreateUser(user); err != nil {
		//return "", err
	}
	key := utils.RandomString(16)
	if _, err := m.IssueToken(username, key); err != nil {
		return "", err
	}
	return key, nil
}
