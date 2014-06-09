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

	auth.HandleFunc("/", MakeAppHandler(authenticate, false)).Methods("GET")
	auth.HandleFunc("/", MakeAppHandler(login, false)).Methods("POST")
	auth.HandleFunc("/", MakeAppHandler(logout, false)).Methods("DELETE")

	photos := r.PathPrefix(fmt.Sprintf("%s/photos",
		settings.Config.ApiPathPrefix)).Subrouter()

	photos.HandleFunc("/", MakeAppHandler(getPhotos, false)).Methods("GET")
	photos.HandleFunc("/", MakeAppHandler(upload, true)).Methods("POST")
	photos.HandleFunc("/{id}", MakeAppHandler(photoDetail, false)).Methods("GET")
	photos.HandleFunc("/{id}", MakeAppHandler(editPhoto, true)).Methods("PUT")
	photos.HandleFunc("/{id}", MakeAppHandler(deletePhoto, true)).Methods("DELETE")

	user := r.PathPrefix(fmt.Sprintf("%s/user",
		settings.Config.ApiPathPrefix)).Subrouter()

	user.HandleFunc("/", MakeAppHandler(signup, false)).Methods("POST")

	r.PathPrefix(settings.Config.PublicPathPrefix).Handler(
		http.FileServer(http.Dir(settings.Config.PublicDir)))

	return session.NewCSRF(r)
}
