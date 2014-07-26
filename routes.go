package photoshare

import (
	"github.com/coopernurse/gorp"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
	"net/http"
	"regexp"
)

var (
	rePhotosByOwnerID = regexp.MustCompile(`^/api/photos/owner/(?P<ownerID>\d+)$`)
	rePhotoDetail     = regexp.MustCompile(`/api/photos/(?P<id>\d+)$`)
	reDeletePhoto     = regexp.MustCompile(`/api/photos/(?P<id>\d+)$`)
	reEditPhotoTitle  = regexp.MustCompile(`/api/photos/(?P<id>\d+)/title$`)
	reEditPhotoTags   = regexp.MustCompile(`/api/photos/(?P<id>\d+)/tags$`)
	reVoteUp          = regexp.MustCompile(`/api/photos/(?P<id>\d+)/upvote$`)
	reVoteDown        = regexp.MustCompile(`/api/photos/(?P<id>\d+)/downvote$`)
	reOwnerFeed       = regexp.MustCompile(`/feeds/owner/(?P<ownerID>\d+)$`)
)

type appHandler func(w http.ResponseWriter, r *request) error

func (h appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTPC(web.C{}, w, r)
}

func (h appHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	req := newRequest(c, r)
	handleError(w, r, h(w, req))
}

type appContext struct {
	config     *appConfig
	ds         *dataStores
	mailer     *mailer
	fs         fileStorage
	sessionMgr sessionManager
	cache      cache
}

// newAppContext creates new AppContext instance
func newAppContext(config *appConfig, dbMap *gorp.DbMap) (*appContext, error) {

	photoDS := newPhotoDataStore(dbMap)
	userDS := newUserDataStore(dbMap)

	ds := &dataStores{
		photos: photoDS,
		users:  userDS,
	}

	fs := newFileStorage(config)
	mailer := newMailer(config)
	cache := newCache(config)

	sessionMgr, err := newSessionManager(config)
	if err != nil {
		return nil, err
	}

	a := &appContext{
		config:     config,
		ds:         ds,
		fs:         fs,
		sessionMgr: sessionMgr,
		mailer:     mailer,
		cache:      cache,
	}
	return a, nil
}

func getRouter(config *appConfig, dbMap *gorp.DbMap) (*web.Mux, error) {

	a, err := newAppContext(config, dbMap)
	if err != nil {
		return nil, err
	}
	r := web.New()

	r.Use(middleware.EnvInit)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AutomaticOptions)

	r.Get("/api/photos/", appHandler(a.getPhotos))
	r.Post("/api/photos/", appHandler(a.upload))
	r.Get("/api/photos/search", appHandler(a.searchPhotos))
	r.Get(rePhotosByOwnerID, appHandler(a.photosByOwnerID))

	r.Get(rePhotoDetail, appHandler(a.photoDetail))
	r.Delete(reDeletePhoto, appHandler(a.deletePhoto))
	r.Patch(reEditPhotoTitle, appHandler(a.editPhotoTitle))
	r.Patch(reEditPhotoTags, appHandler(a.editPhotoTags))
	r.Patch(reVoteUp, appHandler(a.voteUp))
	r.Patch(reVoteDown, appHandler(a.voteDown))

	r.Get("/api/tags/", appHandler(a.getTags))

	r.Get("/api/auth/", appHandler(a.getSessionInfo))
	r.Post("/api/auth/", appHandler(a.login))
	r.Delete("/api/auth/", appHandler(a.logout))
	r.Post("/api/auth/signup", appHandler(a.signup))
	r.Put("/api/auth/recoverpass", appHandler(a.recoverPassword))
	r.Put("/api/auth/changepass", appHandler(a.changePassword))

	r.Get("/feeds/", appHandler(a.latestFeed))
	r.Get("/feeds/popular/", appHandler(a.popularFeed))
	r.Get(reOwnerFeed, appHandler(a.ownerFeed))

	r.Handle("/api/messages/*", messageHandler)
	r.Handle("/*", http.FileServer(http.Dir(config.PublicDir)))
	return r, nil
}
