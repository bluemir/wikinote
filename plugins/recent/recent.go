package recent

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

func (rc *RecentChanges) RegisterRouter(r gin.IRouter) {
	r.GET("/", func(c *gin.Context) {
		lcs := []LastChange{}
		result := rc.db.Order("updated_at desc").Limit(10).Find(&lcs)
		if result.Error != nil {
			c.AbortWithError(http.StatusInternalServerError, result.Error)
			return
		}

		c.JSON(http.StatusOK, lcs)
	})
}
func (rc *RecentChanges) Footer(path string) (template.HTML, error) {
	last := &LastChange{}
	if err := rc.db.Where("path = ?", path).First(last).Error; err != nil {
		return template.HTML(""), err
	}
	return template.HTML("last update: " + last.UpdatedAt.Format(time.RFC3339)), nil
}
