package photoshare

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type params struct {
	vars map[string]string
}

func (p *params) get(name string) string {
	return p.vars[name]
}

func (p *params) getInt(name string) int64 {
	value, _ := strconv.ParseInt(p.vars[name], 10, 0)
	return value
}

type context struct {
	*appConfig
	params *params
	user   *user
}

func (c *context) validate(v validator) error {
	errors := make(map[string]string)
	if err := v.validate(c, errors); err != nil {
		return err
	}
	if len(errors) > 0 {
		return validationFailure{errors}
	}
	return nil
}

func (c *context) getUser(r *http.Request, required bool) (*user, error) {

	if c.user != nil {
		return c.user, nil
	}
	var invalidLogin error

	if required {
		invalidLogin = httpError{http.StatusUnauthorized, "You must be logged in"}
	}

	c.user = &user{}

	userID, err := c.session.readToken(r)
	if err != nil {
		return c.user, err
	}
	if userID == 0 {
		return c.user, invalidLogin
	}
	c.user, err = c.ds.getActiveUser(userID)
	if err != nil {
		if isErrSqlNoRows(err) {
			return c.user, invalidLogin
		}
		return c.user, err
	}
	c.user.IsAuthenticated = true

	return c.user, nil
}

func newContext(cfg *appConfig, r *http.Request) *context {
	c := &context{appConfig: cfg}
	c.params = &params{mux.Vars(r)}
	return c
}
