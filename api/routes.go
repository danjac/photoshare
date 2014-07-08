package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var routeParam = func(r *http.Request, name string) string {
	return mux.Vars(r)[name]
}

var routeParamInt64 = func(r *http.Request, name string) (int64, error) {
	return strconv.ParseInt(routeParam(r, name), 10, 0)
}

func setupRoutes() *mux.Router {

	router := mux.NewRouter()

	api := router.PathPrefix("/api").Subrouter()

	photos := api.PathPrefix("/photos").Subrouter()

	photos.HandleFunc("/", getPhotos).Methods("GET")
	photos.HandleFunc("/", upload).Methods("POST")
	photos.HandleFunc("/search", searchPhotos).Methods("GET")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", photosByOwnerID).Methods("GET")

	photos.HandleFunc("/{id:[0-9]+}", photoDetail).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", deletePhoto).Methods("DELETE")
	photos.HandleFunc("/{id:[0-9]+}/title", editPhotoTitle).Methods("PUT")
	photos.HandleFunc("/{id:[0-9]+}/tags", editPhotoTags).Methods("PUT")
	photos.HandleFunc("/{id:[0-9]+}/upvote", voteUp).Methods("PUT")
	photos.HandleFunc("/{id:[0-9]+}/downvote", voteDown).Methods("PUT")

	api.HandleFunc("/tags/", getTags).Methods("GET")
	api.Handle("/messages/", messageHandler)

	auth := api.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/", authenticate).Methods("GET")
	auth.HandleFunc("/", login).Methods("POST")
	auth.HandleFunc("/", logout).Methods("DELETE")
	auth.HandleFunc("/signup", signup).Methods("POST")
	auth.HandleFunc("/recoverpass", recoverPassword).Methods("PUT")
	auth.HandleFunc("/changepass", changePassword).Methods("PUT")

	feeds := router.PathPrefix("/feeds").Subrouter()

	feeds.HandleFunc("/feeds/", latestFeed).Methods("GET")
	feeds.HandleFunc("/feeds/popular/", popularFeed).Methods("GET")
	feeds.HandleFunc("/feeds/owner/{ownerID:[0-9]+}", ownerFeed).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(config.PublicDir)))

	return router

}
