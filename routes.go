package photoshare

import (
	"github.com/gorilla/mux"
	"net/http"
)

func getRouter(config *appConfig, c *appContext) (*mux.Router, error) {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/").Subrouter()

	photos := api.PathPrefix("/photos/").Subrouter()

	photos.HandleFunc("/", c.handler(getPhotos)).Methods("GET")
	photos.HandleFunc("/", c.handler(upload)).Methods("POST")
	photos.HandleFunc("/search", c.handler(searchPhotos)).Methods("GET")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", c.handler(photosByOwnerID)).Methods("GET")

	photos.HandleFunc("/{id:[0-9]+}", c.handler(getPhotoDetail)).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", c.handler(deletePhoto)).Methods("DELETE")
	photos.HandleFunc("/{id:[0-9]+}/title", c.handler(editPhotoTitle)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/tags", c.handler(editPhotoTags)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/upvote", c.handler(voteUp)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/downvote", c.handler(voteDown)).Methods("PATCH")

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/", c.handler(getSessionInfo)).Methods("GET")
	auth.HandleFunc("/", c.handler(login)).Methods("POST")
	auth.HandleFunc("/", c.handler(logout)).Methods("DELETE")
	auth.HandleFunc("/oauth2/{provider}/url", c.handler(getAuthRedirectURL)).Methods("GET")
	auth.HandleFunc("/oauth2/{provider}/callback/", c.handler(authCallback)).Methods("GET")
	auth.HandleFunc("/signup", c.handler(signup)).Methods("POST")
	auth.HandleFunc("/recoverpass", c.handler(recoverPassword)).Methods("PUT")
	auth.HandleFunc("/changepass", c.handler(changePassword)).Methods("PUT")

	api.HandleFunc("/tags/", c.handler(getTags)).Methods("GET")
	api.Handle("/messages/{path:.*}", messageHandler)

	feeds := r.PathPrefix("/feeds/").Subrouter()

	feeds.HandleFunc("", c.handler(latestFeed)).Methods("GET")
	feeds.HandleFunc("popular/", c.handler(popularFeed)).Methods("GET")
	feeds.HandleFunc("owner/{ownerID:[0-9]+}", c.handler(ownerFeed)).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.PublicDir)))

	return r, nil
}
