package routes

import (
	"github.com/danjac/photoshare/api/models"
)

func signup(c *AppContext) error {

	user := &models.User{}

	if err := c.ParseJSON(user); err != nil {
		return err
	}

	if result, err := user.Validate(); err != nil || !result.OK {
		if err != nil {
			return err
		}
		return c.BadRequest(result)
	}

	if err := user.Insert(); err != nil {
		return err
	}

	if err := c.Login(user); err != nil {
		return err
	}

	return c.OK(user)

}
