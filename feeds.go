package photoshare

import (
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
	"time"
)

func photoFeed(w http.ResponseWriter,
	r *http.Request,
	title string,
	description string,
	link string,
	photos *photoList) error {

	baseURL := baseURL(r)

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

func (a *appContext) latestFeed(_ web.C, w http.ResponseWriter, r *http.Request) error {

	photos, err := a.ds.photos.all(newPage(1), "")

	if err != nil {
		return err
	}

	return photoFeed(w, r, "Latest photos", "Most recent photos", "/latest", photos)
}

func (a *appContext) popularFeed(_ web.C, w http.ResponseWriter, r *http.Request) error {

	photos, err := a.ds.photos.all(newPage(1), "votes")

	if err != nil {
		return err
	}

	return photoFeed(w, r, "Popular photos", "Most upvoted photos", "/popular", photos)
}

func (a *appContext) ownerFeed(c web.C, w http.ResponseWriter, r *http.Request) error {
	ownerID := getIntParam(c, "ownerID")
	owner, err := a.ds.users.getActive(ownerID)
	if err != nil {
		return err
	}

	title := "Feeds for " + owner.Name
	description := "List of feeds for " + owner.Name
	link := fmt.Sprintf("/owner/%d/%s", ownerID, owner.Name)

	photos, err := a.ds.photos.byOwnerID(newPage(1), ownerID)

	if err != nil {
		return err
	}
	return photoFeed(w, r, title, description, link, photos)
}
