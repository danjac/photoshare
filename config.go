package photoshare

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type appConfig struct {
	*settings
	db      *sql.DB
	mailer  *mailer
	ds      dataStore
	fs      fileStorage
	session sessionManager
	auth    authenticator
	cache   cache
}

func newAppConfig() (*appConfig, error) {

	var err error

	settings, err := newSettings()
	if err != nil {
		return nil, err
	}
	config := &appConfig{settings: settings}

	if err := config.initDB(); err != nil {
		return config, err
	}

	config.ds, err = newDataStore(config.db, config.LogSql)
	if err != nil {
		return config, err
	}
	config.fs = newFileStorage(config)
	config.mailer = newMailer(config)
	config.cache = newCache(config)
	config.auth = newAuthenticator(config)

	config.session, err = newSessionManager(config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func (config *appConfig) close() {
	config.db.Close()
}

func (config *appConfig) initDB() error {

	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		config.DBUser,
		config.DBName,
		config.DBPassword,
		config.DBHost,
	))
	if err != nil {
		return err
	}

	config.db = db
	return nil
}

func (config *appConfig) handler(h handlerFunc, loginRequired bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := newContext(config, r)
		if loginRequired {
			_, err := c.getUser(r, true)
			if err != nil {
				handleError(w, r, err)
			}
		}
		handleError(w, r, h(c, w, r))
	}
}

func (config *appConfig) getRouter() http.Handler {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/").Subrouter()

	photos := api.PathPrefix("/photos/").Subrouter()

	photos.HandleFunc("/", config.handler(getPhotos, false)).Methods("GET")
	photos.HandleFunc("/", config.handler(upload, true)).Methods("POST")
	photos.HandleFunc("/search", config.handler(searchPhotos, false)).Methods("GET")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", config.handler(photosByOwnerID, false)).Methods("GET")

	photos.HandleFunc("/{id:[0-9]+}", config.handler(getPhotoDetail, false)).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", config.handler(deletePhoto, true)).Methods("DELETE")
	photos.HandleFunc("/{id:[0-9]+}/title", config.handler(editPhotoTitle, true)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/tags", config.handler(editPhotoTags, true)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/upvote", config.handler(voteUp, true)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/downvote", config.handler(voteDown, true)).Methods("PATCH")

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/", config.handler(getSessionInfo, false)).Methods("GET")
	auth.HandleFunc("/", config.handler(login, false)).Methods("POST")
	auth.HandleFunc("/", config.handler(logout, true)).Methods("DELETE")
	auth.HandleFunc("/oauth2/{provider}/url", config.handler(getAuthRedirectURL, false)).Methods("GET")
	auth.HandleFunc("/oauth2/{provider}/callback/", config.handler(authCallback, false)).Methods("GET")
	auth.HandleFunc("/signup", config.handler(signup, false)).Methods("POST")
	auth.HandleFunc("/recoverpass", config.handler(recoverPassword, false)).Methods("PUT")
	auth.HandleFunc("/changepass", config.handler(changePassword, false)).Methods("PUT")

	api.HandleFunc("/tags/", config.handler(getTags, false)).Methods("GET")
	api.Handle("/messages/{path:.*}", messageHandler)

	feeds := r.PathPrefix("/feeds/").Subrouter()

	feeds.HandleFunc("", config.handler(latestFeed, false)).Methods("GET")
	feeds.HandleFunc("popular/", config.handler(popularFeed, false)).Methods("GET")
	feeds.HandleFunc("owner/{ownerID:[0-9]+}", config.handler(ownerFeed, false)).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.PublicDir)))

	return r

}
