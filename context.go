package photoshare

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// authentication behaviours

type authLevel int

const (
	authLevelIgnore authLevel = iota // we don't need the user in this handler
	authLevelCheck                   // prefetch user, doesn't matter if not logged in
	authLevelLogin                   // user required, 401 if not available
	authLevelAdmin                   // admin required, 401 if no user, 403 if not admin
)

// contains route parameters in a map
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

// request-specific context
// contains the app config so we have access to all the objects we nee
type context struct {
	*app
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

// lazily fetches the current session user
func (ctx *context) authenticate(r *http.Request, level authLevel) (*user, error) {

	if level == authLevelIgnore {
		return &user{}, nil
	}
	var errLoginRequired = httpError{http.StatusUnauthorized, "You must be logged in"}

	var checkAuthLevel = func() error {
		switch level {
		case authLevelLogin:
			if !ctx.user.IsAuthenticated {
				return errLoginRequired
			}
			break
		case authLevelAdmin:
			if !ctx.user.IsAuthenticated {
				return errLoginRequired
			}
			if !ctx.user.IsAdmin {
				return httpError{http.StatusForbidden, "You must be an admin"}
			}
		}
		return nil
	}

	if ctx.user != nil {
		return ctx.user, checkAuthLevel()
	}

	ctx.user = &user{}

	userID, err := ctx.session.readToken(r)
	if err != nil {
		return ctx.user, err
	}
	if userID == 0 {
		return ctx.user, checkAuthLevel()
	}
	ctx.user, err = ctx.datamapper.getActiveUser(userID)
	if err != nil {
		if isErrSqlNoRows(err) {
			return ctx.user, checkAuthLevel()
		}
		return ctx.user, err
	}
	ctx.user.IsAuthenticated = true

	return ctx.user, checkAuthLevel()
}

func newContext(app *app, r *http.Request) *context {
	ctx := &context{app: app}
	ctx.params = &params{mux.Vars(r)}
	return ctx
}
