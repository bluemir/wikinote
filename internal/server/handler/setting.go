package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateCommonSetting(c *gin.Context) {
	/*settings, err := injector.Backends(c).Settings.Get()
	if err != nil {
		c.Error(err)
		return
	}*/
	c.Redirect(http.StatusSeeOther, "/-/admin/common")
}
