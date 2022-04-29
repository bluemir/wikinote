package auth

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

func (m *Manager) IssueToken(username, unhashedKey string, expiredAt *time.Time) (*Token, error) {
	_, ok, err := m.GetUser(username)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.Errorf("user not found")
	}

	token := &Token{
		Username:  username,
		HashedKey: hash(unhashedKey, salt(username)),
		RevokeKey: fmt.Sprintf("%s-%s", xid.New(), hash(username+time.Now().String(), "__revoke__")),
	}
	if err := m.db.Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (m *Manager) RevokeToken(username, unhashedKey string) error {
	user := &User{}
	if err := m.db.Where(&User{Name: username}).Take(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ErrUnauthorized
		}
		return err
	}

	token := &Token{
		Username:  username,
		HashedKey: hash(unhashedKey, user.Salt, m.salt),
	}
	/*
		XXX has bug. see https://github.com/go-gorm/gorm/issues/4879
		if err := m.db.Delete(token).Error; err != nil {
			return err
		}
	*/
	if err := m.db.Model(token).Where(token).Delete(struct {
		UserName  string
		HashedKey string
	}{}).Error; err != nil {
		return err
	}

	return nil
}
func (m *Manager) RevokeTokenAll(username string) error {
	return m.db.Where(&Token{Username: username}).Delete(&Token{Username: username}).Error
}
func (m *Manager) GenerateToken(username string, expiredAt time.Time) (string, error) {
	newKey := hash(xid.New().String(), "__salt__") // TODO Salt

	if _, err := m.IssueToken(username, newKey, &expiredAt); err != nil {
		return "", err
	}
	return newKey, nil
}
