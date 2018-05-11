package recent

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/bluemir/wikinote/plugins"
)

func init() {
	plugins.Register("recent-changes", New)
}

func New(db *gorm.DB, opts map[string]string) plugins.Plugin {

	db.AutoMigrate(&LastChange{})
	logrus.Debugf("init recent-changes, %+v", opts)
	return &RecentChanges{
		db: db,
	}
}

type RecentChanges struct {
	db *gorm.DB
}
type LastChange struct {
	gorm.Model
	Path string
}

func (rc *RecentChanges) AfterWikiSave(path string, data []byte) error {
	lc := &LastChange{
		Path: path,
	}
	result := rc.db.Where("path = ?", path).Assign(&LastChange{
		Model: gorm.Model{
			UpdatedAt: time.Now(),
		},
		Path: path,
	}).FirstOrCreate(lc)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (rc *RecentChanges) RegisterRoute(r *gin.IRoute) {

}
