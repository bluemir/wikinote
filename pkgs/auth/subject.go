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
	if subj.token == nil {
		return "" // just return nil value
	}

	tAttr := &TokenAttr{}
	if err := subj.store.Where(&TokenAttr{
		TokenId: subj.token.ID,
		Key:     key,
	}).Take(tAttr).Error; err != nil {
		if !gorm.IsRecordNotFoundError(err) {
			// TODO log Error
			return ""
		}
	} else {
		// Found
		return tAttr.Value
	}

	user := &User{}
	if err := subj.store.Where(&User{
		Name: subj.token.UserName,
	}).Take(user).Error; err != nil {
		// TODO log Error
		return ""
	}

	uAttr := &UserAttr{}
	if err := subj.store.Where(&UserAttr{
		UserId: user.ID,
		Key:    key,
	}).Take(uAttr).Error; err != nil {
		// TODO log Error
		return ""
	}

	return uAttr.Value
}
