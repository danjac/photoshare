package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/zenazn/goji"
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
	ownerFeedUrl = regexp.MustCompile(`/goji/owner/(?P<ownerID>\d+)$`)
)

func Setup() {

	goji.Get("/api/goji/", getPhotos)
	goji.Get("/api/goji/search", searchPhotos)
	goji.Get(ownerUrl, photosByOwnerID)
	goji.Get(photoUrl, photoDetail)
	goji.Delete(photoUrl, deletePhoto)
	goji.Patch(titleUrl, editPhotoTitle)
	goji.Patch(tagsUrl, editPhotoTags)
	goji.Patch(downvoteUrl, voteDown)
	goji.Patch(upvoteUrl, voteUp)

	goji.Get("/goji/", latestFeed)
	goji.Get("/goji/popular/", popularFeed)
	goji.Get(ownerFeedUrl, ownerFeed)

	goji.Get("/api/auth/", authenticate)
	goji.Post("/api/auth/", login)
	goji.Delete("/api/auth/", logout)

	goji.Post("/api/user/", signup)

	goji.Get("/api/tags/", getTags)

}
