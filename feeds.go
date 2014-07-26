package photoshare

import (
	"fmt"
	"github.com/gorilla/feeds"
	"net/http"
	"strconv"
	"time"
)

func photoFeed(w http.ResponseWriter,
	r *request,
	title string,
	description string,
	link string,
	photos *photoList) error {

	baseURL := r.baseURL()

	feed := &feeds.Feed{
		Title:       title,
		Link:        &feeds.Link{Href: baseURL + link},
		Description: description,
		Created:     time.Now(),
	}

	for _, photo := range photos.Items {

		item := &feeds.Item{
			Id:          strconv.FormatInt(photo.ID, 10),
			Title:       photo.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/#/detail/%d", baseURL, photo.ID)},
			Description: fmt.Sprintf("<img src=\"%s/uploads/thumbnails/%s\">", baseURL, photo.Filename),
			Created:     photo.CreatedAt,
		}
		feed.Add(item)
	}
	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}
	writeBody(w, []byte(atom), http.StatusOK, "application/atom+xml")
	return nil
}

func latestFeed(c *appContext, w http.ResponseWriter, r *request) error {

	photos, err := c.ds.photos.all(newPage(1), "")

	if err != nil {
		return err
	}

	return photoFeed(w, r, "Latest photos", "Most recent photos", "/latest", photos)
}

func popularFeed(c *appContext, w http.ResponseWriter, r *request) error {

	photos, err := c.ds.photos.all(newPage(1), "votes")

	if err != nil {
		return err
	}

	return photoFeed(w, r, "Popular photos", "Most upvoted photos", "/popular", photos)
}

func ownerFeed(c *appContext, w http.ResponseWriter, r *request) error {
	ownerID := r.getIntParam("ownerID")
	owner, err := c.ds.users.getActive(ownerID)
	if err != nil {
		return err
	}

	title := "Feeds for " + owner.Name
	description := "List of feeds for " + owner.Name
	link := fmt.Sprintf("/owner/%d/%s", ownerID, owner.Name)

	photos, err := c.ds.photos.byOwnerID(newPage(1), ownerID)

	if err != nil {
		return err
	}
	return photoFeed(w, r, title, description, link, photos)
}
