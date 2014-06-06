package routes

import (
	"fmt"
	"github.com/danjac/photoshare/api/models"
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

type authHandler func(http.ResponseWriter, *http.Request, *models.User) error

func (fn authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	user, err := session.GetCurrentUser(r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if user == nil {
		render.Status(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	if err := fn(w, r, user); err != nil {
		render.Error(w, r, err)
		return
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
	photos.Handle("/", authHandler(upload)).Methods("POST")
	photos.Handle("/{id}", appHandler(photoDetail)).Methods("GET")
	photos.Handle("/{id}", authHandler(deletePhoto)).Methods("DELETE")

	user := r.PathPrefix(fmt.Sprintf("%s/user",
		settings.Config.ApiPathPrefix)).Subrouter()

	user.Handle("/", appHandler(signup)).Methods("POST")

	r.PathPrefix(settings.Config.PublicPathPrefix).Handler(
		http.FileServer(http.Dir(settings.Config.PublicDir)))

	return session.NewCSRF(r)
}
