package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"regexp"
)

var (
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

func Setup() {

	photos := web.New()
	goji.Handle("/api/photos/*", photos)

	photos.Get("/api/photos/", getPhotos)
	photos.Get("/api/photos/search", searchPhotos)
	photos.Get(ownerUrl, photosByOwnerID)
	photos.Get(photoUrl, photoDetail)
	photos.Delete(photoUrl, deletePhoto)
	photos.Patch(titleUrl, editPhotoTitle)
	photos.Patch(tagsUrl, editPhotoTags)
	photos.Patch(downvoteUrl, voteDown)
	photos.Patch(upvoteUrl, voteUp)

	feeds := web.New()
	goji.Handle("/feeds/*", feeds)

	feeds.Get("/feeds/", latestFeed)
	feeds.Get("/feeds/popular/", popularFeed)
	feeds.Get(ownerFeedUrl, ownerFeed)

	auth := web.New()
	goji.Handle("/api/auth/*", auth)

	auth.Get("/api/auth/", authenticate)
	auth.Post("/api/auth/", login)
	auth.Delete("/api/auth/", logout)

	user := web.New()
	goji.Handle("/api/user/*", auth)

	user.Post("/api/user/", signup)

	tags := web.New()
	goji.Handle("/api/tags/*", tags)

	tags.Get("/api/tags/", getTags)

}
