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
	*config
	params *params
	user   *user
}

func (ctx *context) validate(v validator, r *http.Request) error {
	errors := make(map[string]string)
	if err := v.validate(ctx, r, errors); err != nil {
		return err
	}
	if len(errors) > 0 {
		return validationFailure{errors}
	}
	return nil
}

func (ctx *context) getUser(r *http.Request, required bool) (*user, error) {

	if ctx.user != nil {
		return ctx.user, nil
	}
	var invalidLogin error

	if required {
		invalidLogin = httpError{http.StatusUnauthorized, "You must be logged in"}
	}

	ctx.user = &user{}

	userID, err := ctx.session.readToken(r)
	if err != nil {
		return ctx.user, err
	}
	if userID == 0 {
		return ctx.user, invalidLogin
	}
	ctx.user, err = ctx.datamapper.getActiveUser(userID)
	if err != nil {
		if isErrSqlNoRows(err) {
			return ctx.user, invalidLogin
		}
		return ctx.user, err
	}
	ctx.user.IsAuthenticated = true

	return ctx.user, nil
}

func newContext(cfg *config, r *http.Request) *context {
	ctx := &context{config: cfg}
	ctx.params = &params{mux.Vars(r)}
	return ctx
}
