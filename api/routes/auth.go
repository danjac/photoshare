package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/validation"
)

func logout(c *AppContext) error {

	var (
		token string
		err   error
	)

	if token, err = c.Logout(); err != nil {
		return err
	}

	return c.OK(session.NewSessionInfo(nil, token))

}

func authenticate(c *AppContext) error {

	user, err := c.GetCurrentUser()
	if err != nil {
		return err
	}

	return c.OK(session.NewSessionInfo(user, ""))
}

func login(c *AppContext) error {

	auth := &session.Authenticator{}
	if err := c.ParseJSON(auth); err != nil {
		return err
	}
	user, err := auth.Identify()
	if err != nil {
		if err == session.MissingLoginFields {
			return c.BadRequest("Missing email or password")
		}
		return err
	}
	if user == nil {
		return c.BadRequest("Invalid email or password")
	}

	token, err := c.Login(user)
	if err != nil {
		return err
	}
	return c.OK(session.NewSessionInfo(user, token))
}

func signup(c *AppContext) error {

	user := &models.User{}

	if err := c.ParseJSON(user); err != nil {
		return err
	}

	// ensure nobody tries to make themselves an admin
	user.IsAdmin = false

	validator := &validation.UserValidator{user}

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			return err
		}
		return c.BadRequest(result)
	}

	if err := userMgr.Insert(user); err != nil {
		return err
	}

	token, err := c.Login(user)
	if err != nil {
		return err
	}

	return c.OK(session.NewSessionInfo(user, token))

}
