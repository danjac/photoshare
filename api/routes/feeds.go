package routes

import (
	"fmt"
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/render"
	"github.com/gorilla/feeds"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
	"time"
)

func photoFeed(c web.C, w http.ResponseWriter, r *http.Request, title string, description string, link string, photos []models.Photo) {

	baseURL := baseURL(r)

	feed := &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: baseURL + link},
		Description: description,
		Created:     time.Now(),
	}

	for _, photo := range photos {

		item := &feeds.Item{
			Id:          strconv.FormatInt(photo.ID, 10),
			Title:       photo.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/#/detail/%d", baseURL, photo.ID)},
			Description: fmt.Sprintf("<img src=\"%s/uploads/thumbnails/%s\">", baseURL, photo.Filename),
			Created:     photo.CreatedAt,
		}
		feed.Add(item)
	}

	render.Atom(w, feed, http.StatusOK)

}

func latestFeed(c web.C, w http.ResponseWriter, r *http.Request) {

	list, err := photoMgr.All(1, "")

	if err != nil {
		panic(err)
	}

	photoFeed(c, w, r, "Latest photos", "Most recent photos", "/latest", list.Photos)
}

func popularFeed(c web.C, w http.ResponseWriter, r *http.Request) {

	list, err := photoMgr.All(1, "votes")

	if err != nil {
		panic(err)
	}

	photoFeed(c, w, r, "Popular photos", "Most upvoted photos", "/popular", list.Photos)
}

func ownerFeed(c web.C, w http.ResponseWriter, r *http.Request) {
	ownerID := c.URLParams["ownerID"]
	owner, err := userMgr.GetActive(ownerID)
	if err != nil {
		panic(err)
	}
	if owner == nil {
		http.NotFound(w, r)
		return
	}

	title := "Feeds for " + owner.Name
	description := "List of feeds for " + owner.Name
	link := fmt.Sprintf("/owner/%s/%s", ownerID, owner.Name)

	list, err := photoMgr.ByOwnerID(1, ownerID)

	if err != nil {
		panic(err)
	}
	photoFeed(c, w, r, title, description, link, list.Photos)
}
