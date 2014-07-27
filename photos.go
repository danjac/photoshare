package photoshare

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func deletePhoto(c *appContext, w http.ResponseWriter, r *http.Request, p *params) error {

	user, err := c.getUser(r, true)
	if err != nil {
		return err
	}

	photo, err := c.datastore.photos.get(p.getInt("id"))
	if err != nil {
		return err
	}

	if !photo.canDelete(user) {
		return httpError{http.StatusForbidden, "You're not allowed to delete this photo"}
	}
	if err := c.datastore.photos.remove(photo); err != nil {
		return err
	}

	go func() {
		if err := c.filestore.clean(photo.Filename); err != nil {
			log.Println(err)
		}
	}()

	if err := c.cache.clear(); err != nil {
		return err
	}

	sendMessage(&socketMessage{user.Name, "", photo.ID, "photo_deleted"})
	return renderString(w, http.StatusOK, "Photo deleted")
}

func getPhotoDetail(c *appContext, w http.ResponseWriter, r *http.Request, p *params) error {

	user, err := c.getUser(r, false)
	if err != nil {
		return err
	}

	photo, err := c.datastore.photos.getDetail(p.getInt("id"), user)
	if err != nil {
		return err
	}
	return renderJSON(w, photo, http.StatusOK)

}

func getPhotoToEdit(c *appContext, w http.ResponseWriter, r *http.Request, p *params) (*photo, *user, error) {
	user, err := c.getUser(r, true)
	if err != nil {
		return nil, user, err
	}

	photo, err := c.datastore.photos.get(p.getInt("id"))
	if err != nil {
		return photo, user, err
	}

	if !photo.canEdit(user) {
		return photo, user, httpError{http.StatusForbidden, "You're not allowed to edit this photo"}
	}
	return photo, user, nil
}

func editPhotoTitle(c *appContext, w http.ResponseWriter, r *http.Request, p *params) error {

	photo, user, err := getPhotoToEdit(c, w, r, p)

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

	if err := c.datastore.photos.update(photo); err != nil {
		return err
	}
	sendMessage(&socketMessage{user.Name, "", photo.ID, "photo_updated"})
	return renderString(w, http.StatusOK, "Photo updated")
}

func editPhotoTags(c *appContext, w http.ResponseWriter, r *http.Request, p *params) error {

	photo, user, err := getPhotoToEdit(c, w, r, p)
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

	if err := c.datastore.photos.updateTags(photo); err != nil {
		return err
	}
	sendMessage(&socketMessage{user.Name, "", photo.ID, "photo_updated"})
	return renderString(w, http.StatusOK, "Photo updated")

}

func upload(c *appContext, w http.ResponseWriter, r *http.Request, p *params) error {

	user, err := c.getUser(r, true)
	if err != nil {
		return err
	}

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

	filename, err := c.filestore.store(src, contentType)

	if err != nil {
		if err == errInvalidContentType {
			return httpError{http.StatusBadRequest, err.Error()}
		}
		return err
	}

	photo := &photo{Title: title,
		OwnerID:  user.ID,
		Filename: filename,
		Tags:     tags,
	}

	if err := c.validate(photo); err != nil {
		return err
	}

	if err := c.datastore.photos.create(photo); err != nil {
		return err
	}

	if err := c.cache.clear(); err != nil {
		return err
	}

	sendMessage(&socketMessage{user.Name, "", photo.ID, "photo_uploaded"})
	return renderJSON(w, photo, http.StatusCreated)
}

func searchPhotos(c *appContext, w http.ResponseWriter, r *http.Request, _ *params) error {

	page := getPage(r)
	q := r.FormValue("q")
	qKey := base64.StdEncoding.EncodeToString([]byte(q))
	cacheKey := fmt.Sprintf("photos:search:%s:page:%d", qKey, page.index)

	return c.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := c.datastore.photos.search(page, q)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})

}

func photosByOwnerID(c *appContext, w http.ResponseWriter, r *http.Request, p *params) error {

	page := getPage(r)
	ownerID := p.getInt("ownerID")
	cacheKey := fmt.Sprintf("photos:ownerID:%d:page:%d", ownerID, page.index)

	return c.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := c.datastore.photos.byOwnerID(page, ownerID)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func getPhotos(c *appContext, w http.ResponseWriter, r *http.Request, _ *params) error {

	page := getPage(r)
	orderBy := r.FormValue("orderBy")
	cacheKey := fmt.Sprintf("photos:%s:page:%d", orderBy, page.index)

	return c.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := c.datastore.photos.all(page, orderBy)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func getTags(c *appContext, w http.ResponseWriter, r *http.Request, _ *params) error {
	return c.cache.render(w, http.StatusOK, "tags", func() (interface{}, error) {
		tags, err := c.datastore.photos.getTagCounts()
		if err != nil {
			return tags, err
		}
		return tags, nil
	})

}

func voteDown(c *appContext, w http.ResponseWriter, r *http.Request, p *params) error {
	return vote(c, w, r, p, func(photo *photo) { photo.DownVotes++ })
}

func voteUp(c *appContext, w http.ResponseWriter, r *http.Request, p *params) error {
	return vote(c, w, r, p, func(photo *photo) { photo.UpVotes++ })
}

func vote(c *appContext, w http.ResponseWriter, r *http.Request, p *params, fn func(photo *photo)) error {
	user, err := c.getUser(r, true)
	if err != nil {
		return err
	}

	photo, err := c.datastore.photos.get(p.getInt("id"))
	if err != nil {
		return err
	}

	if !photo.canVote(user) {
		return httpError{http.StatusForbidden, "You're not allowed to vote on this photo"}
	}

	fn(photo)

	if err = c.datastore.photos.update(photo); err != nil {
		return err
	}

	user.registerVote(photo.ID)

	if err = c.datastore.users.update(user); err != nil {
		return err
	}
	return renderString(w, http.StatusOK, "Voting successful")
}
