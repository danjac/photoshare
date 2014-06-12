package session

import (
	"code.google.com/p/xsrftoken"
	"errors"
	"github.com/danjac/photoshare/api/settings"
	"log"
	"net/http"
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

	token, err := cookieReader.Read(r, XsrfCookieName, false)
	if err != nil {
		return false, err
	}

	if token == "" {
		token, err = csrf.Reset(w)
		if err != nil {
			return false, err
		}
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
	return cookieWriter.Write(w, XsrfCookieName, token, false)
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
