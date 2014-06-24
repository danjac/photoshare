package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/validation"
	"strings"
)

func logout(c *Context) *Result {

	if err := c.Logout(); err != nil {
		return c.Error(err)
	}

	return c.OK(session.NewSessionInfo(c.User))

}

func authenticate(c *Context) *Result {

	user, err := c.GetCurrentUser()
	if err != nil {
		return c.Error(err)
	}

	return c.OK(session.NewSessionInfo(user))
}

func login(c *Context) *Result {

	s := &struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{}

	if err := c.ParseJSON(s); err != nil {
		return c.Error(err)
	}

	if s.Identifier == "" || s.Password == "" {
		return c.BadRequest("Missing login details")
	}

	user, err := userMgr.Authenticate(s.Identifier, s.Password)

	if err != nil {
		return c.Error(err)
	}
	if !user.IsAuthenticated {
		return c.BadRequest("Invalid email or password")
	}

	if err := c.Login(user); err != nil {
		return c.Error(err)
	}
	return c.OK(session.NewSessionInfo(user))
}

func signup(c *Context) *Result {

	user := &models.User{}

	if err := c.ParseJSON(user); err != nil {
		return c.Error(err)
	}

	// ensure nobody tries to make themselves an admin
	user.IsAdmin = false

	// email should always be lower case
	user.Email = strings.ToLower(user.Email)

	validator := validation.NewUserValidator(user)

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			return c.Error(err)
		}
		return c.BadRequest(result)
	}

	if err := userMgr.Insert(user); err != nil {
		return c.Error(err)
	}

	if err := c.Login(user); err != nil {
		return c.Error(err)
	}

	user.IsAuthenticated = true

	return c.OK(session.NewSessionInfo(user))

}
