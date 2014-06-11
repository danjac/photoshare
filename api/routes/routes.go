package routes

import (
	"github.com/danjac/photoshare/api/session"
	"github.com/gorilla/mux"
	"net/http"
)

const PublicDir = "./public/"

func GetHandler() http.Handler {

	r := mux.NewRouter()

	auth := r.PathPrefix("/api/auth").Subrouter()

	auth.HandleFunc("/", MakeAppHandler(authenticate, false)).Methods("GET")
	auth.HandleFunc("/", MakeAppHandler(login, false)).Methods("POST")
	auth.HandleFunc("/", MakeAppHandler(logout, false)).Methods("DELETE")

	photos := r.PathPrefix("/api/photos").Subrouter()

	photos.HandleFunc("/", MakeAppHandler(getPhotos, false)).Methods("GET")
	photos.HandleFunc("/", MakeAppHandler(upload, true)).Methods("POST")
	photos.HandleFunc("/{id}", MakeAppHandler(photoDetail, false)).Methods("GET")
	photos.HandleFunc("/{id}", MakeAppHandler(editPhoto, true)).Methods("PUT")
	photos.HandleFunc("/{id}", MakeAppHandler(deletePhoto, true)).Methods("DELETE")

	user := r.PathPrefix("/api/user").Subrouter()

	user.HandleFunc("/", MakeAppHandler(signup, false)).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(PublicDir)))

	return session.NewCSRF(r)
}
