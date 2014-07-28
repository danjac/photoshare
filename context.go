package photoshare

import (
	"github.com/coopernurse/gorp"
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

type handlerFunc func(c *appContext, w http.ResponseWriter, r *http.Request, p *params) error

type appContext struct {
	config  *appConfig
	mailer  *mailer
	ds      dataStore
	fs      fileStorage
	session sessionManager
	auth    authenticator
	cache   cache
}

func (c *appContext) appHandler(h handlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		p := &params{mux.Vars(r)}
		handleError(w, r, h(c, w, r, p))
	}
}

func (c *appContext) validate(v validator) error {
	errors := make(map[string]string)
	if err := v.validate(c, errors); err != nil {
		return err
	}
	if len(errors) > 0 {
		return validationFailure{errors}
	}
	return nil
}

func (c *appContext) getUser(r *http.Request, required bool) (*user, error) {

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

func newAppContext(config *appConfig, dbMap *gorp.DbMap) (*appContext, error) {

	ds := newDataStore(dbMap)
	fs := newFileStorage(config)
	mailer := newMailer(config)
	cache := newCache(config)
	auth := newAuthenticator(config)

	sessionMgr, err := newSessionManager(config)
	if err != nil {
		return nil, err
	}

	c := &appContext{
		config:  config,
		ds:      ds,
		fs:      fs,
		session: sessionMgr,
		mailer:  mailer,
		cache:   cache,
		auth:    auth,
	}
	return c, nil
}
