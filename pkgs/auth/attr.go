package auth

func (m *Manager) SetUserAttr(username string, key, value string) error {
	user := &User{}
	if err := m.store.Where(&User{Name: username}).Take(user).Error; err != nil {
		return err
	}

	if err := m.store.Save(&UserAttr{
		UserId: user.ID,
		Key:    key,
		Value:  value,
	}).Error; err != nil {
		return err
	}
	return nil
}
