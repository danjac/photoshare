package routes

import (
	"encoding/json"
	"fmt"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/settings"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"runtime/debug"
)

func parseJSON(r *http.Request, value interface{}) error {
	return json.NewDecoder(r.Body).Decode(value)
}

func render(w http.ResponseWriter, status int, value interface{}) {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(value)
}

type Recovery struct {
	Handler http.Handler
	Debug   bool
}

func (rec *Recovery) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			f := "PANIC: %s\n%s\n%s"
			log.Printf(f, r, err, stack)
			var msg string
			if rec.Debug {
				msg = fmt.Sprintf(f, r, stack)
			} else {
				msg = "Sorry, an error has occurred"
			}
			http.Error(w, msg, http.StatusInternalServerError)
		}
	}()

	rec.Handler.ServeHTTP(w, r)
}

func NewRecovery(handler http.Handler, debug bool) *Recovery {
	return &Recovery{handler, debug}
}

func Init() http.Handler {

	r := mux.NewRouter()

	auth := r.PathPrefix(fmt.Sprintf("%s/auth",
		settings.Config.ApiPathPrefix)).Subrouter()

	auth.HandleFunc("/", authenticate).Methods("GET")
	auth.HandleFunc("/", login).Methods("POST")
	auth.HandleFunc("/", logout).Methods("DELETE")

	photos := r.PathPrefix(fmt.Sprintf("%s/photos",
		settings.Config.ApiPathPrefix)).Subrouter()

	photos.HandleFunc("/", getPhotos).Methods("GET")
	photos.HandleFunc("/", upload).Methods("POST")
	photos.HandleFunc("/{id}", photoDetail).Methods("GET")
	photos.HandleFunc("/{id}", editPhoto).Methods("PUT")
	photos.HandleFunc("/{id}", deletePhoto).Methods("DELETE")

	user := r.PathPrefix(fmt.Sprintf("%s/user",
		settings.Config.ApiPathPrefix)).Subrouter()

	user.HandleFunc("/", signup).Methods("POST")

	r.PathPrefix(settings.Config.PublicPathPrefix).Handler(
		http.FileServer(http.Dir(settings.Config.PublicDir)))

	return NewRecovery(session.NewCSRF(r), true)
}
