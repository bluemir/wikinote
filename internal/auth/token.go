package auth

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
)

func (m *Manager) IssueToken(username, unhashedKey string) (*Token, error) {
	_, ok, err := m.GetUser(username)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.Errorf("user not found")
	}

	token := &Token{
		UserName:  username,
		HashedKey: hash(unhashedKey, salt(username)),
		RevokeKey: fmt.Sprintf("%s-%s", xid.New(), hash(username+time.Now().String(), "__revoke__")),
	}
	if err := m.db.Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}
func (m *Manager) RevokeToken(revokeKey string) error {
	return m.db.Where(&Token{RevokeKey: revokeKey}).Delete(&Token{RevokeKey: revokeKey}).Error
}
func (m *Manager) RevokeTokenAll(username string) error {
	return m.db.Where(&Token{UserName: username}).Delete(&Token{UserName: username}).Error
}
