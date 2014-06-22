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

	auth.HandleFunc("/", NewAppHandler(authenticate, false)).Methods("GET")
	auth.HandleFunc("/", NewAppHandler(login, false)).Methods("POST")
	auth.HandleFunc("/", NewAppHandler(logout, false)).Methods("DELETE")

	photos := r.PathPrefix("/api/photos").Subrouter()

	photos.HandleFunc("/", NewAppHandler(getPhotos, false)).Methods("GET")
	photos.HandleFunc("/owner/{ownerID}", NewAppHandler(photosByOwnerID, false)).Methods("GET")
	photos.HandleFunc("/search", NewAppHandler(searchPhotos, false)).Methods("GET")
	photos.HandleFunc("/", NewAppHandler(upload, true)).Methods("POST")
	photos.HandleFunc("/{id}", NewAppHandler(photoDetail, false)).Methods("GET")
	photos.HandleFunc("/{id}", NewAppHandler(deletePhoto, true)).Methods("DELETE")
	photos.HandleFunc("/{id}/title", NewAppHandler(editPhotoTitle, true)).Methods("PATCH")
	photos.HandleFunc("/{id}/tags", NewAppHandler(editPhotoTags, true)).Methods("PATCH")
	photos.HandleFunc("/{id}/upvote", NewAppHandler(voteUp, true)).Methods("PATCH")
	photos.HandleFunc("/{id}/downvote", NewAppHandler(voteDown, true)).Methods("PATCH")

	user := r.PathPrefix("/api/user").Subrouter()
	user.HandleFunc("/", NewAppHandler(signup, false)).Methods("POST")

	tags := r.PathPrefix("/api/tags").Subrouter()
	tags.HandleFunc("/", NewAppHandler(getTags, false)).Methods("GET")

	feeds := r.PathPrefix("/feeds").Subrouter()
	feeds.HandleFunc("/", NewAppHandler(latestFeed, false))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(settings.PublicDir)))
	return r
}
