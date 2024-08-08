package auth

func (m *Manager) ListGroups() ([]Group, error) {
	groups := []Group{}
	if err := m.db.Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
