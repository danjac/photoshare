package photoshare

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

func newContext(app *app, r *http.Request, user *user) *context {
	ctx := &context{app: app}
	ctx.params = &params{mux.Vars(r)}
	ctx.user = user
	return ctx
}
