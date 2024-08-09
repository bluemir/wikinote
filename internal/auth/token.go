package auth

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/xid"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Token struct {
	Id        uint   `gorm:"primary_key" json:"id,omitempty"`
	Username  string `json:"username,omitempty"`
	HashedKey string `json:"-,omitempty"`
	RevokeKey string `json:"revoke_key,omitempty"`
	ExpiredAt *time.Time
}
type TokenOpt func(*Token)

func ExpiredAt(t time.Time) func(*Token) {
	return func(token *Token) {
		token.ExpiredAt = &t
	}
}
func ExpiredAfter(d time.Duration) func(*Token) {
	return func(token *Token) {
		t := time.Now().Add(d)
		token.ExpiredAt = &t
	}
}

func (m *Manager) IssueToken(ctx context.Context, username, unhashedKey string, opts ...TokenOpt) (*Token, error) {
	_, ok, err := m.GetUser(ctx, username)
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

	for _, fn := range opts {
		fn(token)
	}
	if err := m.db.WithContext(ctx).Create(token).Error; err != nil {
		return nil, err
	}
	return token, nil
}

func (m *Manager) GenerateToken(ctx context.Context, username string, opts ...TokenOpt) (*Token, string, error) {
	newKey := hash(xid.New().String(), "__salt__") // TODO Salt

	token, err := m.IssueToken(ctx, username, newKey, opts...)
	if err != nil {
		return nil, "", err
	}
	return token, newKey, nil
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
		HashedKey: hash(unhashedKey, user.Salt, m.conf.Salt),
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
