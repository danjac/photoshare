package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type RouteParamsGetter interface {
	String(string) string
	Int(string) (int64, error)
}

var NewRouteParams = func(r *http.Request) RouteParamsGetter {
	return &defaultRouteParams{mux.Vars(r)}
}

type defaultRouteParams struct {
	Params map[string]string
}

func (p *defaultRouteParams) String(name string) string {
	return p.Params[name]
}

func (p *defaultRouteParams) Int(name string) (int64, error) {
	return strconv.ParseInt(p.String(name), 10, 0)
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
	api.PathPrefix("/messages").Handler(messageHandler)

	auth := api.PathPrefix("/auth").Subrouter()

	auth.HandleFunc("/", authenticate).Methods("GET")
	auth.HandleFunc("/", login).Methods("POST")
	auth.HandleFunc("/", logout).Methods("DELETE")
	auth.HandleFunc("/signup", signup).Methods("POST")
	auth.HandleFunc("/recoverpass", recoverPassword).Methods("PUT")
	auth.HandleFunc("/changepass", changePassword).Methods("PUT")

	feeds := router.PathPrefix("/feeds").Subrouter()

	feeds.HandleFunc("/", latestFeed).Methods("GET")
	feeds.HandleFunc("/popular/", popularFeed).Methods("GET")
	feeds.HandleFunc("/owner/{ownerID:[0-9]+}", ownerFeed).Methods("GET")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir(config.PublicDir)))

	return router

}
