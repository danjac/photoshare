package routes

import (
	"fmt"
	"github.com/danjac/photoshare/api/render"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/settings"
	"github.com/gorilla/mux"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if err := fn(w, r); err != nil {
		render.Error(w, r, err)
	}

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
	photos.Handle("/{id}", appHandler(deletePhoto)).Methods("DELETE")

	user := r.PathPrefix(fmt.Sprintf("%s/user",
		settings.Config.ApiPathPrefix)).Subrouter()

	user.Handle("/", appHandler(signup)).Methods("POST")

	r.PathPrefix(settings.Config.PublicPathPrefix).Handler(
		http.FileServer(http.Dir(settings.Config.PublicDir)))

	return session.NewCSRF(r)
}
