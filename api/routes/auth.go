package routes

import (
	"github.com/danjac/photoshare/api/models"
	"net/http"
)

func logout(c *AppContext) {

	if err := c.Logout(); err != nil {
		c.Error(err)
		return
	}

	c.Render(http.StatusOK, "Logged out")

}

func authenticate(c *AppContext) {

	user, err := c.GetCurrentUser()
	if err != nil {
		c.Error(err)
		return
	}
	var status int
	if user == nil {
		status = http.StatusNotFound
	} else {
		status = http.StatusOK
	}

	c.Render(status, user)
}

func login(c *AppContext) {

	auth := &models.Authenticator{}
	if err := c.ParseJSON(auth); err != nil {
		c.Error(err)
		return
	}
	user, err := auth.Identify()
	if err != nil {
		if err == models.MissingLoginFields {
			c.Render(http.StatusBadRequest, "Missing email or password")
			return
		}
		c.Error(err)
		return
	}
	if user == nil {
		c.Render(http.StatusBadRequest, "Invalid email or password")
		return
	}
	if err := c.Login(user); err != nil {
		c.Error(err)
		return
	}
	c.Render(http.StatusOK, user)
}
