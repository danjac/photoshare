package photoshare

import (
	"github.com/coopernurse/gorp"
	"github.com/zenazn/goji/web"
	"net/http"
)

type handlerFunc func(c *appContext, w http.ResponseWriter, r *request) error

type appHandler struct {
	*appContext
	handler       handlerFunc
	loginRequired bool
}

type appContext struct {
	config    *appConfig
	datastore *dataStore
	mailer    *mailer
	filestore fileStorage
	session   sessionManager
	cache     cache
}

func (h *appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTPC(web.C{}, w, r)
}

func (h *appHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	req := newRequest(c, r)
	if h.loginRequired {
		if _, err := h.authenticate(req, true); err != nil {
			handleError(w, r, err)
		}
	}
	handleError(w, r, h.handler(h.appContext, w, req))
}

func (c *appContext) makeAppHandler(h handlerFunc, loginRequired bool) *appHandler {
	return &appHandler{c, h, loginRequired}
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

func (c *appContext) authenticate(r *request, required bool) (*user, error) {

	var invalidLogin error

	if required {
		invalidLogin = httpError{http.StatusUnauthorized, "You must be logged in"}
	}

	if r.user != nil {
		return r.user, nil
	}
	r.user = &user{}

	userID, err := c.session.readToken(r)
	if err != nil {
		return r.user, err
	}
	if userID == 0 {
		return r.user, invalidLogin
	}
	user, err := c.datastore.users.getActive(userID)
	if err != nil {
		if isErrSqlNoRows(err) {
			return r.user, invalidLogin
		}
		return r.user, err
	}
	r.user = user
	r.user.IsAuthenticated = true

	return r.user, nil
}

func newContext(config *appConfig, dbMap *gorp.DbMap) (*appContext, error) {

	ds := newDataStore(dbMap)
	fs := newFileStorage(config)
	mailer := newMailer(config)
	cache := newCache(config)

	sessionMgr, err := newSessionManager(config)
	if err != nil {
		return nil, err
	}

	c := &appContext{
		config:    config,
		datastore: ds,
		filestore: fs,
		session:   sessionMgr,
		mailer:    mailer,
		cache:     cache,
	}
	return c, nil
}
