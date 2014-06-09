package routes

import (
	"github.com/danjac/photoshare/api/models"
	"net/http"
)

func logout(c *AppContext) error {

	if err := c.Logout(); err != nil {
		return err
	}

	return c.OK("Logged out")

}

func authenticate(c *AppContext) error {

	user, err := c.GetCurrentUser()
	if err != nil {
		return err
	}
	var status int
	if user == nil {
		status = http.StatusNotFound
	} else {
		status = http.StatusOK
	}

	return c.Render(status, user)
}

func login(c *AppContext) error {

	auth := &models.Authenticator{}
	if err := c.ParseJSON(auth); err != nil {
		return err
	}
	user, err := auth.Identify()
	if err != nil {
		if err == models.MissingLoginFields {
			return c.BadRequest("Missing email or password")
		}
		return err
	}
	if user == nil {
		return c.BadRequest("Invalid email or password")
	}
	if err := c.Login(user); err != nil {
		return err
	}
	return c.OK(user)
}
