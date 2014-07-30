package photoshare

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

// contains all the objects needed to run the application
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

// our custom handler
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

// the handler should create a new context on each request, and handle any returned
// errors appropriately.
func (cfg *config) handler(h handlerFunc, auth authLevel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleError(w, r, func() error {
			ctx := newContext(cfg, r)
			if _, err := ctx.authenticate(r, auth); err != nil {
				return err
			}
			return h(ctx, w, r)
		}())
	}
}

// generates the routes for the API
func (cfg *config) getRouter() http.Handler {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/").Subrouter()

	photos := api.PathPrefix("/photos/").Subrouter()

	photos.HandleFunc("/", cfg.handler(getPhotos, noAuth)).Methods("GET")
	photos.HandleFunc("/", cfg.handler(upload, userReq)).Methods("POST")
	photos.HandleFunc("/search", cfg.handler(searchPhotos, noAuth)).Methods("GET")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", cfg.handler(photosByOwnerID, noAuth)).Methods("GET")

	photos.HandleFunc("/{id:[0-9]+}", cfg.handler(getPhotoDetail, authReq)).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", cfg.handler(deletePhoto, userReq)).Methods("DELETE")
	photos.HandleFunc("/{id:[0-9]+}/title", cfg.handler(editPhotoTitle, userReq)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/tags", cfg.handler(editPhotoTags, userReq)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/upvote", cfg.handler(voteUp, userReq)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/downvote", cfg.handler(voteDown, userReq)).Methods("PATCH")

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/", cfg.handler(getSessionInfo, authReq)).Methods("GET")
	auth.HandleFunc("/", cfg.handler(login, noAuth)).Methods("POST")
	auth.HandleFunc("/", cfg.handler(logout, userReq)).Methods("DELETE")
	auth.HandleFunc("/oauth2/{provider}/url", cfg.handler(getAuthRedirectURL, noAuth)).Methods("GET")
	auth.HandleFunc("/oauth2/{provider}/callback/", cfg.handler(authCallback, noAuth)).Methods("GET")
	auth.HandleFunc("/signup", cfg.handler(signup, noAuth)).Methods("POST")
	auth.HandleFunc("/recoverpass", cfg.handler(recoverPassword, noAuth)).Methods("PUT")
	auth.HandleFunc("/changepass", cfg.handler(changePassword, noAuth)).Methods("PUT")

	api.HandleFunc("/tags/", cfg.handler(getTags, noAuth)).Methods("GET")
	api.Handle("/messages/{path:.*}", messageHandler)

	feeds := r.PathPrefix("/feeds/").Subrouter()

	feeds.HandleFunc("", cfg.handler(latestFeed, noAuth)).Methods("GET")
	feeds.HandleFunc("popular/", cfg.handler(popularFeed, noAuth)).Methods("GET")
	feeds.HandleFunc("owner/{ownerID:[0-9]+}", cfg.handler(ownerFeed, noAuth)).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(cfg.PublicDir)))

	return r

}
