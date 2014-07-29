package photoshare

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func deletePhoto(c *context, w http.ResponseWriter, r *http.Request) error {

	photo, err := c.ds.getPhoto(c.params.getInt("id"))
	if err != nil {
		return err
	}

	if !photo.canDelete(c.user) {
		return httpError{http.StatusForbidden, "You're not allowed to delete this photo"}
	}
	if err := c.ds.removePhoto(photo); err != nil {
		return err
	}

	go func() {
		if err := c.fs.clean(photo.Filename); err != nil {
			log.Println(err)
		}
	}()

	if err := c.cache.clear(); err != nil {
		return err
	}

	sendMessage(&socketMessage{c.user.Name, "", photo.ID, "photo_deleted"})
	return renderString(w, http.StatusOK, "Photo deleted")
}

func getPhotoDetail(c *context, w http.ResponseWriter, r *http.Request) error {

	user, err := c.getUser(r, false)
	if err != nil {
		return err
	}

	photo, err := c.ds.getPhotoDetail(c.params.getInt("id"), user)
	if err != nil {
		return err
	}
	return renderJSON(w, photo, http.StatusOK)

}

func getPhotoToEdit(c *context, w http.ResponseWriter, r *http.Request) (*photo, error) {

	photo, err := c.ds.getPhoto(c.params.getInt("id"))
	if err != nil {
		return photo, err
	}

	if !photo.canEdit(c.user) {
		return photo, httpError{http.StatusForbidden, "You're not allowed to edit this photo"}
	}
	return photo, nil
}

func editPhotoTitle(c *context, w http.ResponseWriter, r *http.Request) error {

	photo, err := getPhotoToEdit(c, w, r)

	if err != nil {
		return err
	}

	s := &struct {
		Title string `json:"title"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		return err
	}

	photo.Title = s.Title

	if err := c.validate(photo); err != nil {
		return err

	}

	if err := c.ds.updatePhoto(photo); err != nil {
		return err
	}

	sendMessage(&socketMessage{c.user.Name, "", photo.ID, "photo_updated"})
	return renderString(w, http.StatusOK, "Photo updated")
}

func editPhotoTags(c *context, w http.ResponseWriter, r *http.Request) error {

	photo, err := getPhotoToEdit(c, w, r)
	if err != nil {
		return err
	}

	s := &struct {
		Tags []string `json:"tags"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		return err
	}

	photo.Tags = s.Tags
	if err := c.ds.updateTags(photo); err != nil {
		return err
	}

	sendMessage(&socketMessage{c.user.Name, "", photo.ID, "photo_updated"})
	return renderString(w, http.StatusOK, "Photo updated")

}

func upload(c *context, w http.ResponseWriter, r *http.Request) error {

	title := r.FormValue("title")
	taglist := r.FormValue("taglist")
	tags := strings.Split(taglist, " ")

	src, hdr, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			return httpError{http.StatusBadRequest, "Invalid photo"}
		}
		return err
	}
	defer src.Close()

	contentType := hdr.Header["Content-Type"][0]

	filename, err := c.fs.store(src, contentType)

	if err != nil {
		if err == errInvalidContentType {
			return httpError{http.StatusBadRequest, err.Error()}
		}
		return err
	}

	photo := &photo{Title: title,
		OwnerID:  c.user.ID,
		Filename: filename,
		Tags:     tags,
	}

	if err := c.validate(photo); err != nil {
		return err
	}
	if err := c.ds.createPhoto(photo); err != nil {
		return err
	}
	if err := c.cache.clear(); err != nil {
		return err
	}

	sendMessage(&socketMessage{c.user.Name, "", photo.ID, "photo_uploaded"})
	return renderJSON(w, photo, http.StatusCreated)
}

func searchPhotos(c *context, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	q := r.FormValue("q")
	cacheKey := fmt.Sprintf("photos:search:%s:page:%d", q, page.index)

	return c.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := c.ds.searchPhotos(page, q)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})

}

func photosByOwnerID(c *context, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	ownerID := c.params.getInt("ownerID")
	cacheKey := fmt.Sprintf("photos:ownerID:%d:page:%d", ownerID, page.index)

	return c.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := c.ds.getPhotosByOwnerID(page, ownerID)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func getPhotos(c *context, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	orderBy := r.FormValue("orderBy")
	cacheKey := fmt.Sprintf("photos:%s:page:%d", orderBy, page.index)

	return c.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := c.ds.getPhotos(page, orderBy)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func getTags(c *context, w http.ResponseWriter, r *http.Request) error {
	return c.cache.render(w, http.StatusOK, "tags", func() (interface{}, error) {
		tags, err := c.ds.getTagCounts()
		if err != nil {
			return tags, err
		}
		return tags, nil
	})

}

func voteDown(c *context, w http.ResponseWriter, r *http.Request) error {
	return vote(c, w, r, func(photo *photo) { photo.DownVotes++ })
}

func voteUp(c *context, w http.ResponseWriter, r *http.Request) error {
	return vote(c, w, r, func(photo *photo) { photo.UpVotes++ })
}

func vote(c *context, w http.ResponseWriter, r *http.Request, fn func(photo *photo)) error {

	photo, err := c.ds.getPhoto(c.params.getInt("id"))
	if err != nil {
		return err
	}

	if !photo.canVote(c.user) {
		return httpError{http.StatusForbidden, "You're not allowed to vote on this photo"}
	}

	fn(photo)

	c.user.registerVote(photo.ID)

	if err := c.ds.updateMany(photo, c.user); err != nil {
		return err
	}

	return renderString(w, http.StatusOK, "Voting successful")
}
