package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateCommonSetting(c *gin.Context) error {
	/*settings, err := injector.Backends(c).Settings.Get()
	if err != nil {
		return err
	}*/
	c.Redirect(http.StatusSeeOther, "/-/admin/common")
	return nil
}
