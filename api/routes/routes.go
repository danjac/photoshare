package routes

import (
	"fmt"
	"github.com/danjac/photoshare/api/session"
	"github.com/gorilla/mux"
	"net/http"
)

var fileUploadDir string

func Init(uploadsDir, apiPathPrefix, publicPathPrefix, publicDir string) http.Handler {

	fileUploadDir = uploadsDir

	r := mux.NewRouter()

	auth := r.PathPrefix(fmt.Sprintf("%s/auth", apiPathPrefix)).Subrouter()
	auth.HandleFunc("/", authenticate).Methods("GET")
	auth.HandleFunc("/", login).Methods("POST")
	auth.HandleFunc("/", logout).Methods("DELETE")

	photos := r.PathPrefix(fmt.Sprintf("%s/photos", apiPathPrefix)).Subrouter()
	photos.HandleFunc("/", getPhotos).Methods("GET")
	photos.HandleFunc("/", upload).Methods("POST")

	r.PathPrefix(publicPathPrefix).Handler(http.FileServer(http.Dir(publicDir)))

	return session.NewCSRF(r)
}
