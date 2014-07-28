package photoshare

import (
	"github.com/coopernurse/gorp"
	"github.com/gorilla/mux"
	"github.com/juju/errgo"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
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

type authInfo struct {
	name, email string
}

type authenticator interface {
	getRedirectURL(*http.Request, string) (string, error)
	getUserInfo(*http.Request, string) (*authInfo, error)
}

func newAuthenticator(config *appConfig) authenticator {
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	a := &defaultAuthenticator{config}
	return a
}

type defaultAuthenticator struct {
	config *appConfig
}

func (a *defaultAuthenticator) getAuthProvider(r *http.Request, providerName string) (common.Provider, error) {
	gomniauth.WithProviders(
		google.New(a.config.GoogleAuthKey,
			a.config.GoogleAuthSecret,
			getBaseURL(r)+"/api/auth/oauth2/google/callback/",
		),
	)
	provider, err := gomniauth.Provider(providerName)
	if err != nil {
		return provider, errgo.Mask(err)
	}
	return provider, nil
}

func (a *defaultAuthenticator) getRedirectURL(r *http.Request, providerName string) (string, error) {
	provider, err := a.getAuthProvider(r, providerName)
	if err != nil {
		return "", errgo.Mask(err)
	}
	state := gomniauth.NewState("after", "success")
	url, err := provider.GetBeginAuthURL(state, nil)
	if err != nil {
		return url, errgo.Mask(err)
	}
	return url, nil
}

func (a *defaultAuthenticator) getUserInfo(r *http.Request, providerName string) (*authInfo, error) {
	provider, err := a.getAuthProvider(r, providerName)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	m := make(objx.Map)
	if r.Form == nil {
		r.ParseForm()
	}
	for k, v := range r.Form {
		m.Set(k, v)
	}
	creds, err := provider.CompleteAuth(m)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	user, err := provider.GetUser(creds)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	info := &authInfo{
		name:  user.Name(),
		email: user.Email(),
	}

	return info, nil
}

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
