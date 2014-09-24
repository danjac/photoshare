package photoshare

import (
	"database/sql"
	"github.com/gorilla/mux"
	"net/http"
)

// authentication behaviours

type authLevel int

const (
	authLevelIgnore authLevel = iota // we don't need the user in this handler
	authLevelCheck                   // prefetch user, doesn't matter if not logged in
	authLevelLogin                   // user required, 401 if not available
	authLevelAdmin                   // admin required, 401 if no user, 403 if not admin
)

// contains all the objects needed to run the application
type app struct {
	cfg        *config
	db         *sql.DB
	mailer     *mailer
	router     *mux.Router
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

	app.initRouter()

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
			user, err := app.authenticate(r, level)
			if err != nil {
				return err
			}
			return h(newContext(app, r, user), w, r)
		}())
	}
}

// lazily fetches the current session user
func (app *app) authenticate(r *http.Request, level authLevel) (*user, error) {

	if level == authLevelIgnore {
		return &user{}, nil
	}
	var errLoginRequired = httpError{http.StatusUnauthorized, "You must be logged in"}

	var checkAuthLevel = func(user *user) error {
		switch level {
		case authLevelLogin:
			if !user.IsAuthenticated {
				return errLoginRequired
			}
			break
		case authLevelAdmin:
			if !user.IsAuthenticated {
				return errLoginRequired
			}
			if !user.IsAdmin {
				return httpError{http.StatusForbidden, "You must be an admin"}
			}
		}
		return nil
	}

	user := &user{}

	userID, err := app.session.readToken(r)
	if err != nil {
		return user, err
	}
	if userID == 0 {
		return user, checkAuthLevel(user)
	}
	user, err = app.datamapper.getActiveUser(userID)
	if err != nil {
		if isErrSqlNoRows(err) {
			return user, checkAuthLevel(user)
		}
		return nil, err
	}
	user.IsAuthenticated = true

	return user, checkAuthLevel(user)
}

// generates the routes for the API
func (app *app) initRouter() {

	app.router = mux.NewRouter()

	api := app.router.PathPrefix("/api/").Subrouter()

	photos := api.PathPrefix("/photos/").Subrouter()

	photos.HandleFunc("/", app.handler(getPhotos, authLevelIgnore)).Methods("GET").Name("photos")
	photos.HandleFunc("/", app.handler(upload, authLevelLogin)).Methods("POST").Name("photos")
	photos.HandleFunc("/search", app.handler(searchPhotos, authLevelIgnore)).Methods("GET").Name("search")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", app.handler(photosByOwnerID, authLevelIgnore)).Methods("GET").Name("owner")

	photos.HandleFunc("/{id:[0-9]+}", app.handler(getPhotoDetail, authLevelCheck)).Methods("GET").Name("photoDetail")
	photos.HandleFunc("/{id:[0-9]+}", app.handler(deletePhoto, authLevelLogin)).Methods("DELETE").Name("deletePhoto")
	photos.HandleFunc("/{id:[0-9]+}/title", app.handler(editPhotoTitle, authLevelLogin)).Methods("PATCH").Name("editPhotoTitle")
	photos.HandleFunc("/{id:[0-9]+}/tags", app.handler(editPhotoTags, authLevelLogin)).Methods("PATCH").Name("editPhotoTags")
	photos.HandleFunc("/{id:[0-9]+}/upvote", app.handler(voteUp, authLevelLogin)).Methods("PATCH").Name("upvote")
	photos.HandleFunc("/{id:[0-9]+}/downvote", app.handler(voteDown, authLevelLogin)).Methods("PATCH").Name("downvote")

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/", app.handler(getSessionInfo, authLevelCheck)).Methods("GET").Name("sessionInfo")
	auth.HandleFunc("/", app.handler(login, authLevelIgnore)).Methods("POST").Name("login")
	auth.HandleFunc("/", app.handler(logout, authLevelLogin)).Methods("DELETE").Name("logout")
	auth.HandleFunc("/signup", app.handler(signup, authLevelIgnore)).Methods("POST").Name("signup")
	auth.HandleFunc("/recoverpass", app.handler(recoverPassword, authLevelIgnore)).Methods("PUT").Name("recoverPassword")
	auth.HandleFunc("/changepass", app.handler(changePassword, authLevelIgnore)).Methods("PUT").Name("changePassword")

	auth.HandleFunc("/oauth2/{provider}/url", app.handler(getAuthRedirectURL, authLevelIgnore)).Methods("GET")
	auth.HandleFunc("/oauth2/{provider}/callback/", app.handler(authCallback, authLevelIgnore)).Methods("GET")

	api.HandleFunc("/tags/", app.handler(getTags, authLevelIgnore)).Methods("GET").Name("tags")
	api.Handle("/messages/{path:.*}", messageHandler).Name("messages")

	feeds := app.router.PathPrefix("/feeds/").Subrouter()

	feeds.HandleFunc("", app.handler(latestFeed, authLevelIgnore)).Methods("GET").Name("latestFeed")
	feeds.HandleFunc("popular/", app.handler(popularFeed, authLevelIgnore)).Methods("GET").Name("popularFeed")
	feeds.HandleFunc("owner/{ownerID:[0-9]+}", app.handler(ownerFeed, authLevelIgnore)).Methods("GET").Name("ownerFeed")

	app.router.PathPrefix("/").Handler(http.FileServer(http.Dir(app.cfg.PublicDir)))

}
