package photoshare

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

// contains all the objects needed to run the application
type app struct {
	cfg        *config
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

func newApp() (*app, error) {

	var err error

	app := &app{}

	app.cfg, err = newConfig()
	if err != nil {
		return app, err
	}

	if err := app.initDB(); err != nil {
		return app, err
	}

	app.datamapper, err = newDataMapper(app.db, app.cfg.LogSql)
	if err != nil {
		return app, err
	}
	app.filestore = newFileStorage(app.cfg)
	app.mailer = newMailer(app.cfg)
	app.cache = newCache(app.cfg)
	app.auth = newAuthenticator(app.cfg)

	app.session, err = newSessionManager(app.cfg)
	if err != nil {
		return app, err
	}

	return app, nil
}

func (app *app) close() {
	app.db.Close()
}

func (app *app) initDB() error {

	db, err := dbConnect(app.cfg.DBUser,
		app.cfg.DBPassword,
		app.cfg.DBName,
		app.cfg.DBHost)
	if err != nil {
		return err
	}
	app.db = db
	return nil
}

// the handler should create a new context on each request, and handle any returned
// errors appropriately.
func (app *app) handler(h handlerFunc, level authLevel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleError(w, r, func() error {
			ctx := newContext(app, r)
			if _, err := ctx.authenticate(r, level); err != nil {
				return err
			}
			return h(ctx, w, r)
		}())
	}
}

// generates the routes for the API
func (app *app) getRouter() http.Handler {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/").Subrouter()

	photos := api.PathPrefix("/photos/").Subrouter()

	photos.HandleFunc("/", app.handler(getPhotos, noAuth)).Methods("GET")
	photos.HandleFunc("/", app.handler(upload, userReq)).Methods("POST")
	photos.HandleFunc("/search", app.handler(searchPhotos, noAuth)).Methods("GET")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", app.handler(photosByOwnerID, noAuth)).Methods("GET")

	photos.HandleFunc("/{id:[0-9]+}", app.handler(getPhotoDetail, authReq)).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", app.handler(deletePhoto, userReq)).Methods("DELETE")
	photos.HandleFunc("/{id:[0-9]+}/title", app.handler(editPhotoTitle, userReq)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/tags", app.handler(editPhotoTags, userReq)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/upvote", app.handler(voteUp, userReq)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/downvote", app.handler(voteDown, userReq)).Methods("PATCH")

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/", app.handler(getSessionInfo, authReq)).Methods("GET")
	auth.HandleFunc("/", app.handler(login, noAuth)).Methods("POST")
	auth.HandleFunc("/", app.handler(logout, userReq)).Methods("DELETE")
	auth.HandleFunc("/oauth2/{provider}/url", app.handler(getAuthRedirectURL, noAuth)).Methods("GET")
	auth.HandleFunc("/oauth2/{provider}/callback/", app.handler(authCallback, noAuth)).Methods("GET")
	auth.HandleFunc("/signup", app.handler(signup, noAuth)).Methods("POST")
	auth.HandleFunc("/recoverpass", app.handler(recoverPassword, noAuth)).Methods("PUT")
	auth.HandleFunc("/changepass", app.handler(changePassword, noAuth)).Methods("PUT")

	api.HandleFunc("/tags/", app.handler(getTags, noAuth)).Methods("GET")
	api.Handle("/messages/{path:.*}", messageHandler)

	feeds := r.PathPrefix("/feeds/").Subrouter()

	feeds.HandleFunc("", app.handler(latestFeed, noAuth)).Methods("GET")
	feeds.HandleFunc("popular/", app.handler(popularFeed, noAuth)).Methods("GET")
	feeds.HandleFunc("owner/{ownerID:[0-9]+}", app.handler(ownerFeed, noAuth)).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(app.cfg.PublicDir)))

	return r

}
