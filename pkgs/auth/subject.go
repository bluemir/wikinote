package auth

import (
	"github.com/jinzhu/gorm"
)

type subject struct {
	*Manager
	token *Token
}

func (m *Manager) Subject(token *Token) Subject {
	// token may be nil for guest
	return &subject{m, token}
}

func (subj *subject) Attr(key string) string {
	attr := &Attr{}

	if subj.token == nil {
		return "" // just return nil value
	}

	if err := subj.store.Where(&TokenAttr{
		TokenId: subj.token.ID,
		Key:     key,
	}).Take(attr).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		// TODO log Error
		return ""
	}

	user := &User{}
	if err := subj.store.Where(&User{
		Name: subj.token.UserName,
	}).Take(user).Error; err != nil {
		// TODO log Error
		return ""
	}
	if err := subj.store.Where(&UserAttr{
		UserId: user.ID,
		Key:    key,
	}).Take(attr).Error; err != nil {
		// TODO log Error
		return ""
	}

	return attr.Value
}
