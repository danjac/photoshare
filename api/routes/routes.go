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

	auth := r.PathPrefix(fmt.Sprintf("%s/auth", settings.Config.ApiPathPrefix)).Subrouter()
	auth.HandleFunc("/", authenticate).Methods("GET")
	auth.HandleFunc("/", login).Methods("POST")
	auth.HandleFunc("/", logout).Methods("DELETE")

	photos := r.PathPrefix(fmt.Sprintf("%s/photos", settings.Config.ApiPathPrefix)).Subrouter()
	photos.HandleFunc("/", getPhotos).Methods("GET")
	photos.HandleFunc("/", upload).Methods("POST")

	r.PathPrefix(settings.Config.PublicPathPrefix).Handler(http.FileServer(http.Dir(settings.Config.PublicDir)))

	return session.NewCSRF(r)
}
