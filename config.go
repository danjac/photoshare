package photoshare

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/coopernurse/gorp"
	"github.com/danryan/env"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path"
)

type settings struct {
	DBName     string `env:"key=DB_NAME required=true"`
	DBUser     string `env:"key=DB_USER required=true"`
	DBPassword string `env:"key=DB_PASS required=true"`
	DBHost     string `env:"key=DB_HOST default=localhost"`

	TestDBName     string `env:"key=TEST_DB_NAME"`
	TestDBUser     string `env:"key=TEST_DB_USER"`
	TestDBPassword string `env:"key=TEST_DB_PASS"`
	TestDBHost     string `env:"key=TEST_DB_HOST"`

	LogSql bool `env:"key=LOG_SQL default=false"`

	SmtpName          string `env:"key=SMTP_NAME"`
	SmtpPassword      string `env:"key=SMTP_PASS"`
	SmtpUser          string `env:"key=SMTP_USER"`
	SmtpHost          string `env:"key=SMTP_HOST default=localhost"`
	SmtpPort          int    `env:"key=SMTP_PORT default=25"`
	SmtpDefaultSender string `env:"key=DEFAULT_EMAIL_SENDER default=webmaster@localhost"`

	BaseDir       string `env:"key=BASE_DIR"`
	PublicDir     string `env:"key=PUBLIC_DIR"`
	UploadsDir    string `env:"key=UPLOADS_DIR"`
	ThumbnailsDir string `env:"key=THUMBNAILS_DIR"`
	TemplatesDir  string `env:"key=TEMPLATES_DIR"`

	PrivateKey string `env:"key=PRIVATE_KEY required=true"`
	PublicKey  string `env:"key=PUBLIC_KEY required=true"`

	MemcacheHost string `env:"key=MEMCACHE_HOST default=0.0.0.0:11211"`

	GoogleAuthKey    string `env:"key=GOOGLE_AUTH_KEY"`
	GoogleAuthSecret string `env:"key=GOOGLE_AUTH_SECRET"`

	ServerPort int `env:"key=PORT default=5000"`
}

type appConfig struct {
	*settings
	dbMap   *gorp.DbMap
	mailer  *mailer
	ds      dataStore
	fs      fileStorage
	session sessionManager
	auth    authenticator
	cache   cache
}

func newAppConfig() (*appConfig, error) {

	settings := &settings{}

	if err := env.Process(settings); err != nil {
		return nil, err
	}

	config := &appConfig{settings: settings}

	if config.TestDBName == "" {
		config.TestDBName = config.DBName + "_test"
	}

	if config.TestDBUser == "" {
		config.TestDBUser = config.DBUser
	}

	if config.TestDBPassword == "" {
		config.TestDBPassword = config.DBPassword
	}

	if config.TestDBHost == "" {
		config.TestDBHost = config.DBHost
	}

	if config.TestDBName == config.DBName {
		return config, errors.New("test DB name same as DB name")
	}

	if config.BaseDir == "" {
		config.BaseDir = getDefaultBaseDir()
	}

	if config.PublicDir == "" {
		config.PublicDir = path.Join(config.BaseDir, "public")
	}

	if config.UploadsDir == "" {
		config.UploadsDir = path.Join(config.PublicDir, "uploads")
	}

	if config.ThumbnailsDir == "" {
		config.ThumbnailsDir = path.Join(config.UploadsDir, "thumbnails")
	}

	if config.TemplatesDir == "" {
		config.TemplatesDir = path.Join(config.BaseDir, "templates")
	}

	if err := config.initDB(); err != nil {
		return config, err
	}

	if err := config.initDB(); err != nil {
		return config, err
	}

	config.ds = newDataStore(config.dbMap)
	config.fs = newFileStorage(config)
	config.mailer = newMailer(config)
	config.cache = newCache(config)
	config.auth = newAuthenticator(config)

	config.session, _ = newSessionManager(config)

	return config, nil
}

func (config *appConfig) close() {
	config.dbMap.Db.Close()
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

	config.dbMap, err = initDB(db, config.LogSql)
	if err != nil {
		return err
	}
	return nil
}

func (config *appConfig) handler(h handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleError(w, r, h(config.makeContext(r), w, r))
	}
}

func (config *appConfig) makeContext(r *http.Request) *context {

	c := &context{appConfig: config}
	c.params = &params{mux.Vars(r)}
	return c
}

func (config *appConfig) getRouter() http.Handler {

	r := mux.NewRouter()

	api := r.PathPrefix("/api/").Subrouter()

	photos := api.PathPrefix("/photos/").Subrouter()

	photos.HandleFunc("/", config.handler(getPhotos)).Methods("GET")
	photos.HandleFunc("/", config.handler(upload)).Methods("POST")
	photos.HandleFunc("/search", config.handler(searchPhotos)).Methods("GET")
	photos.HandleFunc("/owner/{ownerID:[0-9]+}", config.handler(photosByOwnerID)).Methods("GET")

	photos.HandleFunc("/{id:[0-9]+}", config.handler(getPhotoDetail)).Methods("GET")
	photos.HandleFunc("/{id:[0-9]+}", config.handler(deletePhoto)).Methods("DELETE")
	photos.HandleFunc("/{id:[0-9]+}/title", config.handler(editPhotoTitle)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/tags", config.handler(editPhotoTags)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/upvote", config.handler(voteUp)).Methods("PATCH")
	photos.HandleFunc("/{id:[0-9]+}/downvote", config.handler(voteDown)).Methods("PATCH")

	auth := api.PathPrefix("/auth/").Subrouter()

	auth.HandleFunc("/", config.handler(getSessionInfo)).Methods("GET")
	auth.HandleFunc("/", config.handler(login)).Methods("POST")
	auth.HandleFunc("/", config.handler(logout)).Methods("DELETE")
	auth.HandleFunc("/oauth2/{provider}/url", config.handler(getAuthRedirectURL)).Methods("GET")
	auth.HandleFunc("/oauth2/{provider}/callback/", config.handler(authCallback)).Methods("GET")
	auth.HandleFunc("/signup", config.handler(signup)).Methods("POST")
	auth.HandleFunc("/recoverpass", config.handler(recoverPassword)).Methods("PUT")
	auth.HandleFunc("/changepass", config.handler(changePassword)).Methods("PUT")

	api.HandleFunc("/tags/", config.handler(getTags)).Methods("GET")
	api.Handle("/messages/{path:.*}", messageHandler)

	feeds := r.PathPrefix("/feeds/").Subrouter()

	feeds.HandleFunc("", config.handler(latestFeed)).Methods("GET")
	feeds.HandleFunc("popular/", config.handler(popularFeed)).Methods("GET")
	feeds.HandleFunc("owner/{ownerID:[0-9]+}", config.handler(ownerFeed)).Methods("GET")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.PublicDir)))

	return r

}

func getDefaultBaseDir() string {
	defaultBaseDir, err := os.Getwd()
	if err != nil {
		defaultBaseDir = "."
	}
	return defaultBaseDir
}
