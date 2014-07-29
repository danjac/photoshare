package photoshare

import (
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

type handlerFunc func(c *context, w http.ResponseWriter, r *http.Request) error

type context struct {
	*appConfig
	params *params
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

	var invalidLogin error

	if required {
		invalidLogin = httpError{http.StatusUnauthorized, "You must be logged in"}
	}

	user := &user{}

	userID, err := c.session.readToken(r)
	if err != nil {
		return user, err
	}
	if userID == 0 {
		return user, invalidLogin
	}
	user, err = c.ds.getActiveUser(userID)
	if err != nil {
		if isErrSqlNoRows(err) {
			return user, invalidLogin
		}
		return user, err
	}
	user.IsAuthenticated = true

	return user, nil
}
