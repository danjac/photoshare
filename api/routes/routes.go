package routes

import (
	"fmt"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/settings"
	"github.com/gorilla/mux"
	"net/http"
)

func Init() http.Handler {

	r := mux.NewRouter()

	auth := r.PathPrefix(fmt.Sprintf("%s/auth",
		settings.Config.ApiPathPrefix)).Subrouter()

	auth.HandleFunc("/", NewAppHandler(authenticate)).Methods("GET")
	auth.HandleFunc("/", NewAppHandler(login)).Methods("POST")
	auth.HandleFunc("/", NewAppHandler(logout)).Methods("DELETE")

	photos := r.PathPrefix(fmt.Sprintf("%s/photos",
		settings.Config.ApiPathPrefix)).Subrouter()

	photos.HandleFunc("/", NewAppHandler(getPhotos)).Methods("GET")
	photos.HandleFunc("/", NewAppHandler(upload)).Methods("POST")
	photos.HandleFunc("/{id}", NewAppHandler(photoDetail)).Methods("GET")
	photos.HandleFunc("/{id}", NewAppHandler(editPhoto)).Methods("PUT")
	photos.HandleFunc("/{id}", NewAppHandler(deletePhoto)).Methods("DELETE")

	user := r.PathPrefix(fmt.Sprintf("%s/user",
		settings.Config.ApiPathPrefix)).Subrouter()

	user.HandleFunc("/", NewAppHandler(signup)).Methods("POST")

	r.PathPrefix(settings.Config.PublicPathPrefix).Handler(
		http.FileServer(http.Dir(settings.Config.PublicDir)))

	return session.NewCSRF(r)
}
