package photoshare

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type config struct {
	*settings
	db         *sql.DB
	mailer     *mailer
	datamapper dataMapper
	filestore  fileStorage
	session    sessionManager
	auth       authenticator
	cache      cache
}

type handlerFunc func(c *context, w http.ResponseWriter, r *http.Request) error

func newConfig() (*config, error) {

	var err error

	settings, err := newSettings()
	if err != nil {
		return nil, err
	}
	cfg := &config{settings: settings}

	if err := cfg.initDB(); err != nil {
		return cfg, err
	}

	cfg.datamapper, err = newDataMapper(cfg.db, cfg.LogSql)
	if err != nil {
		return cfg, err
	}
	cfg.filestore = newFileStorage(cfg)
	cfg.mailer = newMailer(cfg)
	cfg.cache = newCache(cfg)
	cfg.auth = newAuthenticator(cfg)

	cfg.session, err = newSessionManager(cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func (cfg *config) close() {
	cfg.db.Close()
}

func (cfg *config) initDB() error {

	db, err := dbConnect(cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBHost)
	if err != nil {
		return err
	}
	cfg.db = db
	return nil
}

func (cfg *config) handler(h handlerFunc, loginRequired bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleError(w, r, func() error {
			ctx := newContext(cfg, r)
			if loginRequired {
				_, err := ctx.getUser(r, true)
				if err != nil {
					return err
				}
			}
			return h(ctx, w, r)
		}())
	}
}

func (cfg *config) getRouter() http.Handler {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/").Subrouter()

	photos := api.PathPrefix("/photos/").Subrouter()

	photos.HandleFunc("/", cfg.handler(getPhotos, false)).Methods("GET")
	photos.HandleFunc("/", cfg.handler(upload, true)).Methods("POST")
	photos.HandleFunc("/search", cfg.handler(searchPhotos, false)).Methods("GET")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", cfg.handler(photosByOwnerID, false)).Methods("GET")

	photos.HandleFunc("/{id:[0-9]+}", cfg.handler(getPhotoDetail, false)).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", cfg.handler(deletePhoto, true)).Methods("DELETE")
	photos.HandleFunc("/{id:[0-9]+}/title", cfg.handler(editPhotoTitle, true)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/tags", cfg.handler(editPhotoTags, true)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/upvote", cfg.handler(voteUp, true)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/downvote", cfg.handler(voteDown, true)).Methods("PATCH")

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/", cfg.handler(getSessionInfo, false)).Methods("GET")
	auth.HandleFunc("/", cfg.handler(login, false)).Methods("POST")
	auth.HandleFunc("/", cfg.handler(logout, true)).Methods("DELETE")
	auth.HandleFunc("/oauth2/{provider}/url", cfg.handler(getAuthRedirectURL, false)).Methods("GET")
	auth.HandleFunc("/oauth2/{provider}/callback/", cfg.handler(authCallback, false)).Methods("GET")
	auth.HandleFunc("/signup", cfg.handler(signup, false)).Methods("POST")
	auth.HandleFunc("/recoverpass", cfg.handler(recoverPassword, false)).Methods("PUT")
	auth.HandleFunc("/changepass", cfg.handler(changePassword, false)).Methods("PUT")

	api.HandleFunc("/tags/", cfg.handler(getTags, false)).Methods("GET")
	api.Handle("/messages/{path:.*}", messageHandler)

	feeds := r.PathPrefix("/feeds/").Subrouter()

	feeds.HandleFunc("", cfg.handler(latestFeed, false)).Methods("GET")
	feeds.HandleFunc("popular/", cfg.handler(popularFeed, false)).Methods("GET")
	feeds.HandleFunc("owner/{ownerID:[0-9]+}", cfg.handler(ownerFeed, false)).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.PublicDir)))

	return r

}

func dbConnect(user, pwd, name, host string) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s",
		user,
		name,
		pwd,
		host,
	))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
