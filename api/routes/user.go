package routes

import (
	"github.com/danjac/photoshare/api/models"
)

func signup(c *AppContext) {

	user := &models.User{}

	if err := c.ParseJSON(user); err != nil {
		c.Error(err)
		return
	}

	if result, err := user.Validate(); err != nil || !result.OK {
		if err != nil {
			c.Error(err)
			return
		}
		c.BadRequest(result)
	}

	if err := user.Insert(); err != nil {
		c.Error(err)
		return
	}

	if err := c.Login(user); err != nil {
		c.Error(err)
		return
	}

	c.OK(user)

}
