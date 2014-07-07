package routes

import (
	"github.com/danjac/photoshare/api/email"
	"github.com/danjac/photoshare/api/models"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"regexp"
)

var (
	mailer       = email.NewMailer()
	photoMgr     = models.NewPhotoManager()
	userMgr      = models.NewUserManager()
	ownerUrl     = regexp.MustCompile(`/api/photos/owner/(?P<ownerID>\d+)$`)
	photoUrl     = regexp.MustCompile(`/api/photos/(?P<id>\d+)$`)
	titleUrl     = regexp.MustCompile(`/api/photos/(?P<id>\d+)/title$`)
	tagsUrl      = regexp.MustCompile(`/api/photos/(?P<id>\d+)/tags$`)
	downvoteUrl  = regexp.MustCompile(`/api/photos/(?P<id>\d+)/downvote$`)
	upvoteUrl    = regexp.MustCompile(`/api/photos/(?P<id>\d+)/upvote$`)
	ownerFeedUrl = regexp.MustCompile(`/feeds/owner/(?P<ownerID>\d+)$`)
)

func init() {

	photos := web.New()
	photos.Get("/api/photos/", getPhotos)
	photos.Get("/api/photos/search", searchPhotos)
	photos.Get(ownerUrl, photosByOwnerID)
	photos.Get(photoUrl, photoDetail)

	photos.Post("/api/photos/", upload)
	photos.Delete(photoUrl, deletePhoto)
	photos.Patch(titleUrl, editPhotoTitle)
	photos.Patch(tagsUrl, editPhotoTags)
	photos.Patch(downvoteUrl, voteDown)
	photos.Patch(upvoteUrl, voteUp)

	goji.Handle("/api/photos/*", photos)

	tags := web.New()
	tags.Get("/api/tags/", getTags)

	goji.Handle("/api/tags/*", tags)

	messages := web.New()
	messages.Handle("/api/messages/*", messageHandler)

	goji.Handle("/api/messages/*", messages)

	auth := web.New()
	auth.Get("/api/auth/", authenticate)
	auth.Post("/api/auth/", login)
	auth.Delete("/api/auth/", logout)
	auth.Post("/api/auth/signup", signup)
	auth.Put("/api/auth/recoverpass", recoverPassword)
	auth.Put("/api/auth/changepass", changePassword)

	goji.Handle("/api/auth/*", auth)

	feeds := web.New()
	feeds.Get("/feeds/", latestFeed)
	feeds.Get("/feeds/popular/", popularFeed)
	feeds.Get(ownerFeedUrl, ownerFeed)

	goji.Handle("/feeds/*", feeds)

}
