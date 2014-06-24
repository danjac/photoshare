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

	photos.HandleFunc("/", NewAppHandler(upload, true)).Methods("POST")
	photos.HandleFunc("/", NewAppHandler(getPhotos, false)).Methods("GET")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", NewAppHandler(photosByOwnerID, false)).Methods("GET")
	photos.HandleFunc("/search", NewAppHandler(searchPhotos, false)).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", NewAppHandler(photoDetail, false)).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", NewAppHandler(deletePhoto, true)).Methods("DELETE")
	photos.HandleFunc("/{id:[0-9]+}/title", NewAppHandler(editPhotoTitle, true)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/tags", NewAppHandler(editPhotoTags, true)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/upvote", NewAppHandler(voteUp, true)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/downvote", NewAppHandler(voteDown, true)).Methods("PATCH")

	user := r.PathPrefix("/api/user").Subrouter()

	user.HandleFunc("/", NewAppHandler(signup, false)).Methods("POST")

	tags := r.PathPrefix("/api/tags").Subrouter()

	tags.HandleFunc("/", NewAppHandler(getTags, false)).Methods("GET")

	feeds := r.PathPrefix("/feeds").Subrouter()

	feeds.HandleFunc("/", NewAppHandler(latestFeed, false))
	feeds.HandleFunc("/owner/{ownerID:[0-9]+}", NewAppHandler(ownerFeed, false)).Methods("GET")
	feeds.HandleFunc("/popular", NewAppHandler(popularFeed, false))

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(settings.PublicDir)))
	return r
}
