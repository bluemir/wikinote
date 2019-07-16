package auth

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
	}).Take(attr).Error; err != nil {
		return ""
	}

	return attr.Value
}
