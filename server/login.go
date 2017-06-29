package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "/login.html", Data(c))
}
func HandleLogin(c *gin.Context) {
	id := c.PostForm("id")
	password := c.PostForm("password")

	user, err := Backend(c).User().Get(id)
	if err != nil {
		FlashMessage(c).Warn("Error on get user: %s", err.Error())
		c.Redirect(http.StatusSeeOther, "/!/auth/login")
		return
	}

	if !user.Password.Check(password) {
		FlashMessage(c).Warn("wrong password")
		c.Redirect(http.StatusSeeOther, "/!/auth/login")
		return
	}
	// login success
	Session(c).Login(user)

	FlashMessage(c).Info("Welcome %s", user.Id)

	c.Redirect(http.StatusSeeOther, "/")
}
func HandleLogout(c *gin.Context) {
	Session(c).Logout()

	FlashMessage(c).Info("Bye(logout successfully)")
	c.Redirect(http.StatusSeeOther, "/")
}
