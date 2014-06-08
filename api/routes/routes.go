package routes

import (
	"encoding/json"
	"fmt"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/settings"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if err := fn(w, r); err != nil {
		log.Println(err, r)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func parseJSON(r *http.Request, value interface{}) error {
	return json.NewDecoder(r.Body).Decode(value)
}

func render(w http.ResponseWriter, status int, value interface{}) error {
	w.WriteHeader(status)
	w.Header().Set("Content-type", "application/json")
	return json.NewEncoder(w).Encode(value)
}

func Init() http.Handler {

	r := mux.NewRouter()

	auth := r.PathPrefix(fmt.Sprintf("%s/auth",
		settings.Config.ApiPathPrefix)).Subrouter()

	auth.Handle("/", appHandler(authenticate)).Methods("GET")
	auth.Handle("/", appHandler(login)).Methods("POST")
	auth.Handle("/", appHandler(logout)).Methods("DELETE")

	photos := r.PathPrefix(fmt.Sprintf("%s/photos",
		settings.Config.ApiPathPrefix)).Subrouter()

	photos.Handle("/", appHandler(getPhotos)).Methods("GET")
	photos.Handle("/", appHandler(upload)).Methods("POST")
	photos.Handle("/{id}", appHandler(photoDetail)).Methods("GET")
	photos.Handle("/{id}", appHandler(editPhoto)).Methods("PUT")
	photos.Handle("/{id}", appHandler(deletePhoto)).Methods("DELETE")

	user := r.PathPrefix(fmt.Sprintf("%s/user",
		settings.Config.ApiPathPrefix)).Subrouter()

	user.Handle("/", appHandler(signup)).Methods("POST")

	r.PathPrefix(settings.Config.PublicPathPrefix).Handler(
		http.FileServer(http.Dir(settings.Config.PublicDir)))

	return session.NewCSRF(r)
}
