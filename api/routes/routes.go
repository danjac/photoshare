package routes

import (
	"github.com/danjac/photoshare/api/session"
	"github.com/gorilla/mux"
	"net/http"
)

func Init() http.Handler {
	r := mux.NewRouter()

	auth := r.PathPrefix("/api/auth").Subrouter()
	auth.HandleFunc("/", authenticate).Methods("GET")
	auth.HandleFunc("/", login).Methods("POST")
	auth.HandleFunc("/", logout).Methods("DELETE")

	photos := r.PathPrefix("/api/photos").Subrouter()
	photos.HandleFunc("/", getPhotos).Methods("GET")
	photos.HandleFunc("/", upload).Methods("POST")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	return session.NewCSRF(r)
}
