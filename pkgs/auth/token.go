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
	if err := m.store.Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}
func (m *Manager) ListToken(username string) ([]Token, error) {
	result := []Token{}
	err := m.store.Where(&Token{UserName: username}).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, err
}
func (m *Manager) RevokeToken(revokeKey string) error {
	err := m.store.Where(&Token{RevokeKey: revokeKey}).Delete(&Token{}).Error
	if err != nil {
		return err
	}
	return nil
}
