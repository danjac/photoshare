package routes

import (
	"fmt"
	"github.com/gorilla/feeds"
	"net/http"
	"strconv"
	"time"
)

func latestFeed(c *Context) *Result {

	var scheme string
	if c.Request.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, c.Request.Host)

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

	atom, err := feed.ToAtom()
	if err != nil {
		return c.Error(err)
	}

	c.Response.WriteHeader(http.StatusOK)
	c.Response.Header().Set("Content-Type", "application/atom+xml")
	c.Response.Write([]byte(atom))
	return nil

}
