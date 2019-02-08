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

func New(opts map[string]string, store plugins.FileAttrStore) plugins.Plugin {
	logrus.Debugf("init recent-changes, %+v", opts)
	return &RecentChanges{
		store: store,
	}
}

type RecentChanges struct {
	store plugins.FileAttrStore
}
type LastChange struct {
	gorm.Model
	Path string
}

func (rc *RecentChanges) OnPostSave(path string, data []byte, store plugins.FileAttr) error {
	// for last update use utc (sortable)
	return store.Set("plugin.wikinote.bluemir.me/last-change", time.Now().UTC().Format(time.RFC3339))
}

func (rc *RecentChanges) RegisterRouter(r gin.IRouter) {
	r.GET("/", func(c *gin.Context) {
		attrs, err := rc.store.Where(&plugins.FindOptions{
			Namespace: "plugin.wikinote.bluemir.me",
			Key:       "last-change",
		}).SortBy(plugins.VALUE, plugins.ASC).Limit(10).Find()
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
		}

		c.JSON(http.StatusOK, attrs)
	})
}

func (rc *RecentChanges) Footer(path string, store plugins.FileAttr) (template.HTML, error) {
	tStr, err := store.Get("plugin.wikinote.bluemir.me/last-change")
	if err != nil {
		return template.HTML(""), err
	}

	t, err := time.Parse(time.RFC3339, tStr)
	if err != nil {
		return template.HTML(""), err
	}

	return template.HTML("last update: " + t.Local().Format(time.RFC3339)), nil
}
