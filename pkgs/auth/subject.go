package auth

import (
	"github.com/jinzhu/gorm"
)

type subject struct {
	*Manager
	token *Token
}

func (m *Manager) Subject(token *Token) Subject {
	return &subject{m, token}
}

func (subj *subject) Attr(key string) string {
	attr := &Attr{}
	if err := subj.store.Where(&Attr{
		ID:  subj.token.ID,
		Key: key,
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
	if err := subj.store.Where(&Attr{
		ID:  user.ID,
		Key: key,
	}).Take(attr).Error; err != nil {
		// TODO log Error
		return ""
	}

	return attr.Value
}
