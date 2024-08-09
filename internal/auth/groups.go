package auth

import "context"

func (m *Manager) ListGroups(ctx context.Context) ([]Group, error) {
	groups := []Group{}
	if err := m.db.WithContext(ctx).Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}
