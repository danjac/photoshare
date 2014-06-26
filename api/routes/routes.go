package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/zenazn/goji"
	"gopkg.in/igm/sockjs-go.v2/sockjs"
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

	goji.Get("/api/photos/", getPhotos)
	goji.Post("/api/photos/", upload)
	goji.Get("/api/photos/search", searchPhotos)
	goji.Get(ownerUrl, photosByOwnerID)
	goji.Get(photoUrl, photoDetail)
	goji.Delete(photoUrl, deletePhoto)

	goji.Patch(titleUrl, editPhotoTitle)
	goji.Patch(tagsUrl, editPhotoTags)
	goji.Patch(downvoteUrl, voteDown)
	goji.Patch(upvoteUrl, voteUp)

	goji.Get("/api/auth/", authenticate)
	goji.Post("/api/auth/", login)
	goji.Delete("/api/auth/", logout)
	goji.Post("/api/user/", signup)

	goji.Get("/api/tags/", getTags)

	goji.Get("/feeds/", latestFeed)
	goji.Get("/feeds/popular/", popularFeed)
	goji.Get(ownerFeedUrl, ownerFeed)

	goji.Handle("/api/messages/*", sockjs.NewHandler("/api/messages",
		sockjs.DefaultOptions,
		messageHandler))

}
