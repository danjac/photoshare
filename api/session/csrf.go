package session

import (
	"code.google.com/p/xsrftoken"
	"github.com/danjac/photoshare/api/render"
	"net/http"
	"time"
)

const (
	XsrfCookieName = "csrf_token"
	XsrfHeaderName = "X-CSRF-Token"
)

type CSRF struct {
	SecretKey string
	Handler   http.Handler
}

func NewCSRF(handler http.Handler) *CSRF {
	csrf := &CSRF{string(hashKey), handler}
	return csrf
}

func (csrf *CSRF) Validate(w http.ResponseWriter, r *http.Request) bool {

	var token string

	cookie, err := r.Cookie(XsrfCookieName)
	if err != nil || cookie.Value == "" {
		token = csrf.Reset(w)
	} else {
		token = cookie.Value
	}

	if r.Method == "GET" || r.Method == "OPTIONS" || r.Method == "HEAD" {
		return true
	}

	return token != "" && token == r.Header.Get(XsrfHeaderName)
}

func (csrf *CSRF) Reset(w http.ResponseWriter) string {
	token := csrf.Generate()
	csrf.Save(w, token)
	return token
}

func (csrf *CSRF) Generate() string {
	return xsrftoken.Generate(csrf.SecretKey, "xsrf", "POST")
}

func (csrf *CSRF) Save(w http.ResponseWriter, token string) {

	expires := time.Now().AddDate(0, 0, 1)

	cookie := &http.Cookie{
		Name:       XsrfCookieName,
		Value:      token,
		Path:       "/",
		MaxAge:     86400,
		Expires:    expires,
		RawExpires: expires.Format(time.UnixDate),
	}

	http.SetCookie(w, cookie)
}

func (csrf *CSRF) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !csrf.Validate(w, r) {
		render.Ping(w, http.StatusForbidden, "Invalid CSRF header")
		return
	}
	csrf.Handler.ServeHTTP(w, r)
}
