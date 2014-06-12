package session

import (
	"code.google.com/p/xsrftoken"
	"errors"
	"github.com/danjac/photoshare/api/settings"
	"log"
	"net/http"
	"time"
)

const (
	XsrfCookieName = "csrf_token"
	XsrfHeaderName = "X-CSRF-Token"
)

var InvalidCSRFToken = errors.New("Invalid CSRF token")

type CSRF struct {
	SecretKey string
	Handler   http.Handler
}

func NewCSRF(handler http.Handler) *CSRF {
	csrf := &CSRF{settings.HashKey, handler}
	return csrf
}

func (csrf *CSRF) Validate(w http.ResponseWriter, r *http.Request) (bool, error) {

	var token string

	cookie, err := r.Cookie(XsrfCookieName)

	if err != nil || cookie.Value == "" {
		if err := sCookie.Decode(XsrfCookieName, cookie.Value, &token); err != nil {
			return false, nil
		}
		token, err = csrf.Reset(w)
		if err != nil {
			return false, err
		}
	} else {
		token = cookie.Value
	}

	if r.Method == "GET" || r.Method == "OPTIONS" || r.Method == "HEAD" {
		return true, nil
	}

	return token != "" && token == r.Header.Get(XsrfHeaderName), nil
}

func (csrf *CSRF) Reset(w http.ResponseWriter) (string, error) {
	token := csrf.Generate()
	if err := csrf.Save(w, token); err != nil {
		return token, err
	}
	return token, nil
}

func (csrf *CSRF) Generate() string {
	return xsrftoken.Generate(csrf.SecretKey, "xsrf", "POST")
}

func (csrf *CSRF) Save(w http.ResponseWriter, token string) error {

	expires := time.Now().AddDate(0, 0, 1)

	encoded, err := sCookie.Encode(XsrfCookieName, token)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:       XsrfCookieName,
		Value:      encoded,
		Path:       "/",
		MaxAge:     86400,
		Expires:    expires,
		RawExpires: expires.Format(time.UnixDate),
	}

	http.SetCookie(w, cookie)
	return nil
}

func (csrf *CSRF) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if ok, err := csrf.Validate(w, r); !ok || err != nil {
		if err != nil {
			log.Println("ERROR:", err)
		}
		http.Error(w, InvalidCSRFToken.Error(), http.StatusForbidden)
		return
	}
	csrf.Handler.ServeHTTP(w, r)
}
