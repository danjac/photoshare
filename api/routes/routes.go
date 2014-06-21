package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/settings"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	photoMgr = models.NewPhotoManager()
	userMgr  = models.NewUserManager()
)

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
	photos.HandleFunc("/{id}", MakeAppHandler(deletePhoto, true)).Methods("DELETE")
	photos.HandleFunc("/{id}/title", MakeAppHandler(editPhotoTitle, true)).Methods("PATCH")
	photos.HandleFunc("/{id}/tags", MakeAppHandler(editPhotoTags, true)).Methods("PATCH")
	photos.HandleFunc("/{id}/upvote", MakeAppHandler(voteUp, true)).Methods("PATCH")
	photos.HandleFunc("/{id}/downvote", MakeAppHandler(voteDown, true)).Methods("PATCH")

	user := r.PathPrefix("/api/user").Subrouter()
	user.HandleFunc("/", MakeAppHandler(signup, false)).Methods("POST")

	tags := r.PathPrefix("/api/tags").Subrouter()
	tags.HandleFunc("/", MakeAppHandler(getTags, false)).Methods("GET")

	feeds := r.PathPrefix("/feeds").Subrouter()
	feeds.HandleFunc("/", MakeAppHandler(latestFeed, false))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(settings.PublicDir)))
	return r
}
