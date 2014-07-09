package api

import (
	"github.com/zenazn/goji"
	"net/http"
)

func initRoutes() {

	goji.Get("/api/photos/", getPhotos)
	goji.Post("/api/photos/", upload)
	goji.Get("/api/photos/search", searchPhotos)
	goji.Get("/api/photos/owner/:ownerID", photosByOwnerID)

	goji.Get("/api/photos/:id", photoDetail)
	goji.Delete("/api/photos/:id", deletePhoto)
	goji.Patch("/api/photos/:id/title", editPhotoTitle)
	goji.Patch("/api/photos/:id/tags", editPhotoTags)
	goji.Patch("/api/photos/:id/upvote", voteUp)
	goji.Patch("/api/photos/:id/downvote", voteDown)

	goji.Get("/api/tags/", getTags)
	goji.Handle("/api/messages/*", messageHandler)

	goji.Get("/api/auth/", authenticate)
	goji.Post("/api/auth/", login)
	goji.Delete("/api/auth/", logout)
	goji.Post("/api/auth/signup", signup)
	goji.Put("/api/auth/recoverpass", recoverPassword)
	goji.Put("/api/auth/changepass", changePassword)

	goji.Get("/feeds/", latestFeed)
	goji.Get("/feeds/popular/", popularFeed)
	goji.Get("/feeds/owner/:ownerID", ownerFeed)

	goji.Handle("/*", http.FileServer(http.Dir(config.PublicDir)))

}
