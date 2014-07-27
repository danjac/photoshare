package photoshare

import (
	"github.com/gorilla/mux"
	"net/http"
)

func getRouter(config *appConfig, c *appContext) (*mux.Router, error) {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/").Subrouter()

	photos := api.PathPrefix("/photos/").Subrouter()

	photos.HandleFunc("/", c.appHandler(getPhotos)).Methods("GET")
	photos.HandleFunc("/", c.appHandler(upload)).Methods("POST")
	photos.HandleFunc("/search", c.appHandler(searchPhotos)).Methods("GET")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", c.appHandler(photosByOwnerID)).Methods("GET")

	photos.HandleFunc("/{id:[0-9]+}", c.appHandler(getPhotoDetail)).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", c.appHandler(deletePhoto)).Methods("DELETE")
	photos.HandleFunc("/{id:[0-9]+}/title", c.appHandler(editPhotoTitle)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/tags", c.appHandler(editPhotoTags)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/upvote", c.appHandler(voteUp)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/downvote", c.appHandler(voteDown)).Methods("PATCH")

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/", c.appHandler(getSessionInfo)).Methods("GET")
	auth.HandleFunc("/", c.appHandler(login)).Methods("POST")
	auth.HandleFunc("/", c.appHandler(logout)).Methods("DELETE")
	auth.HandleFunc("/signup", c.appHandler(signup)).Methods("POST")
	auth.HandleFunc("/recoverpass", c.appHandler(recoverPassword)).Methods("PUT")
	auth.HandleFunc("/changepass", c.appHandler(changePassword)).Methods("PUT")

	api.HandleFunc("/tags/", c.appHandler(getTags)).Methods("GET")
	api.Handle("/messages/{path:.*}", messageHandler)

	feeds := r.PathPrefix("/feeds/").Subrouter()

	feeds.HandleFunc("", c.appHandler(latestFeed)).Methods("GET")
	feeds.HandleFunc("popular/", c.appHandler(popularFeed)).Methods("GET")
	feeds.HandleFunc("owner/{ownerID:[0-9]+}", c.appHandler(ownerFeed)).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.PublicDir)))

	return r, nil
}
