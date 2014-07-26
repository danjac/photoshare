package photoshare

import (
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

func getRouter(config *appConfig, c *context) (*web.Mux, error) {

	r := web.New()

	r.Use(middleware.EnvInit)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AutomaticOptions)

	r.Get("/api/photos/", c.makeAppHandler(getPhotos, false))
	r.Post("/api/photos/", c.makeAppHandler(upload, true))
	r.Get("/api/photos/search", c.makeAppHandler(searchPhotos, false))
	r.Get(rePhotosByOwnerID, c.makeAppHandler(photosByOwnerID, false))

	r.Get(rePhotoDetail, c.makeAppHandler(getPhotoDetail, false))
	r.Delete(reDeletePhoto, c.makeAppHandler(deletePhoto, true))
	r.Patch(reEditPhotoTitle, c.makeAppHandler(editPhotoTitle, true))
	r.Patch(reEditPhotoTags, c.makeAppHandler(editPhotoTags, true))
	r.Patch(reVoteUp, c.makeAppHandler(voteUp, true))
	r.Patch(reVoteDown, c.makeAppHandler(voteDown, true))

	r.Get("/api/tags/", c.makeAppHandler(getTags, false))

	r.Get("/api/auth/", c.makeAppHandler(getSessionInfo, false))
	r.Post("/api/auth/", c.makeAppHandler(login, false))
	r.Delete("/api/auth/", c.makeAppHandler(logout, true))
	r.Post("/api/auth/signup", c.makeAppHandler(signup, false))
	r.Put("/api/auth/recoverpass", c.makeAppHandler(recoverPassword, false))
	r.Put("/api/auth/changepass", c.makeAppHandler(changePassword, false))

	r.Get("/feeds/", c.makeAppHandler(latestFeed, false))
	r.Get("/feeds/popular/", c.makeAppHandler(popularFeed, false))
	r.Get(reOwnerFeed, c.makeAppHandler(ownerFeed, false))

	r.Handle("/api/messages/*", messageHandler)
	r.Handle("/*", http.FileServer(http.Dir(config.PublicDir)))
	return r, nil
}
