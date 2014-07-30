package photoshare

import (
	"github.com/juju/errgo"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/common"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	"github.com/stretchr/signature"
	"net/http"
)

type authInfo struct {
	name, email string
}

type authenticator interface {
	getRedirectURL(*http.Request, string) (string, error)
	getUserInfo(*http.Request, string) (*authInfo, error)
}

func newAuthenticator(cfg *config) authenticator {
	gomniauth.SetSecurityKey(signature.RandomKey(64))
	a := &defaultAuthenticator{cfg}
	return a
}

type defaultAuthenticator struct {
	cfg *config
}

func (a *defaultAuthenticator) getAuthProvider(r *http.Request, providerName string) (common.Provider, error) {
	gomniauth.WithProviders(
		google.New(a.cfg.GoogleClientID,
			a.cfg.GoogleSecret,
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
