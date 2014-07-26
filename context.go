package photoshare

import (
	"github.com/coopernurse/gorp"
	"github.com/zenazn/goji/web"
	"net/http"
)

type handlerFunc func(c *context, w http.ResponseWriter, r *request) error

type appHandler struct {
	*context
	handler       handlerFunc
	loginRequired bool
}

type context struct {
	config     *appConfig
	ds         *dataStores
	mailer     *mailer
	fs         fileStorage
	sessionMgr sessionManager
	cache      cache
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
	handleError(w, r, h.handler(h.context, w, req))
}

func (c *context) makeAppHandler(h handlerFunc, loginRequired bool) *appHandler {
	return &appHandler{c, h, loginRequired}
}

func (c *context) authenticate(r *request, required bool) (*user, error) {

	var invalidLogin error

	if required {
		invalidLogin = httpError{http.StatusUnauthorized, "You must be logged in"}
	}

	if r.user != nil {
		return r.user, nil
	}
	r.user = &user{}

	userID, err := c.sessionMgr.readToken(r)
	if err != nil {
		return r.user, err
	}
	if userID == 0 {
		return r.user, invalidLogin
	}
	user, err := c.ds.users.getActive(userID)
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

func newContext(config *appConfig, dbMap *gorp.DbMap) (*context, error) {

	photoDS := newPhotoDataStore(dbMap)
	userDS := newUserDataStore(dbMap)

	ds := &dataStores{
		photos: photoDS,
		users:  userDS,
	}

	fs := newFileStorage(config)
	mailer := newMailer(config)
	cache := newCache(config)

	sessionMgr, err := newSessionManager(config)
	if err != nil {
		return nil, err
	}

	c := &context{
		config:     config,
		ds:         ds,
		fs:         fs,
		sessionMgr: sessionMgr,
		mailer:     mailer,
		cache:      cache,
	}
	return c, nil
}
