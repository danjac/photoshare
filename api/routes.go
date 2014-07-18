package api

import (
	"database/sql"
	"fmt"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
)

// return HttpError{http.StatusForbidden}
type HttpError struct {
	Status int
}

func (h HttpError) Error() string {
	return fmt.Sprintf("HTTP %d:%s", h.Status, http.StatusText(h.Status))
}

type AppHandler func(c web.C, w http.ResponseWriter, r *http.Request) error

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h(web.C{}, w, r)
	handleError(w, r, err)
}

func (h AppHandler) ServeHTTPC(c web.C, w http.ResponseWriter, r *http.Request) {
	err := h(c, w, r)
	handleError(w, r, err)
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	httpError, ok := err.(HttpError)
	if ok {
		http.Error(w, http.StatusText(httpError.Status), httpError.Status)
		return
	}

	var status int

	switch err {
	case sql.ErrNoRows:
		status = http.StatusNotFound
	default:
		status = http.StatusInternalServerError
	}
	log.Println(err) // more sophisticated logging needed
	http.Error(w, http.StatusText(status), status)
}

func initRoutes() {

	goji.Get("/api/photos/", AppHandler(getPhotos))
	goji.Post("/api/photos/", AppHandler(upload))
	goji.Get("/api/photos/search", AppHandler(searchPhotos))
	goji.Get("/api/photos/owner/:ownerID", AppHandler(photosByOwnerID))

	goji.Get("/api/photos/:id", AppHandler(photoDetail))
	goji.Delete("/api/photos/:id", AppHandler(deletePhoto))
	goji.Patch("/api/photos/:id/title", AppHandler(editPhotoTitle))
	goji.Patch("/api/photos/:id/tags", AppHandler(editPhotoTags))
	goji.Patch("/api/photos/:id/upvote", AppHandler(voteUp))
	goji.Patch("/api/photos/:id/downvote", AppHandler(voteDown))

	goji.Get("/api/tags/", AppHandler(getTags))

	goji.Get("/api/auth/", AppHandler(authenticate))
	goji.Post("/api/auth/", AppHandler(login))
	goji.Delete("/api/auth/", AppHandler(logout))
	goji.Post("/api/auth/signup", AppHandler(signup))
	goji.Put("/api/auth/recoverpass", AppHandler(recoverPassword))
	goji.Put("/api/auth/changepass", AppHandler(changePassword))

	goji.Get("/feeds/", AppHandler(latestFeed))
	goji.Get("/feeds/popular/", AppHandler(popularFeed))
	goji.Get("/feeds/owner/:ownerID", AppHandler(ownerFeed))

	goji.Handle("/api/messages/*", messageHandler)
	goji.Handle("/*", http.FileServer(http.Dir(config.PublicDir)))

}
