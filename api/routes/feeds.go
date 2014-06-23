package routes

import (
	"fmt"
	"github.com/danjac/photoshare/api/models"
	"github.com/gorilla/feeds"
	"strconv"
	"time"
)

func photoFeed(c *Context, title string, description string, link string, photos []models.Photo) *Result {

	baseURL := c.BaseURL()

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

	return c.Atom(feed)

}

func latestFeed(c *Context) *Result {

	list, err := photoMgr.All(1, "")

	if err != nil {
		return c.Error(err)
	}

	return photoFeed(c, "Latest photos", "Most recent photos", "/latest", list.Photos)
}

func popularFeed(c *Context) *Result {

	list, err := photoMgr.All(1, "votes")

	if err != nil {
		return c.Error(err)
	}

	return photoFeed(c, "Popular photos", "Most upvoted photos", "/popular", list.Photos)
}

func ownerFeed(c *Context) *Result {
	ownerID := c.Param("ownerID")
	owner, err := userMgr.GetActive(ownerID)
	if err != nil {
		return c.Error(err)
	}
	if owner == nil {
		return c.NotFound("No user found")
	}

	title := "Feeds for " + owner.Name
	description := "List of feeds for " + owner.Name
	link := fmt.Sprintf("/owner/%s/%s", ownerID, owner.Name)

	list, err := photoMgr.ByOwnerID(1, ownerID)

	if err != nil {
		return c.Error(err)
	}
	return photoFeed(c, title, description, link, list.Photos)
}
