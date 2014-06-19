package routes

import (
	"fmt"
	"github.com/gorilla/feeds"
	"strconv"
	"time"
)

func latestFeed(c *Context) *Result {

	baseURL := c.BaseURL()

	photos, err := photoMgr.All(1)

	if err != nil {
		return c.Error(err)
	}
	feed := &feeds.Feed{
		Title:       "Latest photos",
		Link:        &feeds.Link{Href: baseURL + "/"},
		Description: "Latest photos",
		Created:     time.Now(),
	}

	for _, photo := range photos {

		item := &feeds.Item{
			Id:          strconv.FormatInt(photo.ID, 10),
			Title:       photo.Title,
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/#/detail/%d", baseURL, photo.ID)},
			Description: fmt.Sprintf("<img src=\"%s/uploads/thumbnails/%s\">", baseURL, photo.Photo),
			Created:     photo.CreatedAt,
		}
		feed.Add(item)
	}

	return c.Atomize(feed)

}
