package file

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MetaInfo struct {
	gorm.Model
	Path string
	// TODO will be added
}

func (m *manager) saveTime(path string) error {
	meta := &MetaInfo{
		Path: path,
	}
	result := m.db.Where("path = ?", path).Assign(&MetaInfo{
		Model: gorm.Model{
			UpdatedAt: time.Now(),
		},
		Path: path,
	}).FirstOrCreate(meta)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (m *manger) RecentChanges(n int) ([]MetaInfo, error) {
	return nil, nil
}
