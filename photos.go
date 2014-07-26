package photoshare

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func deletePhoto(c *appContext, w http.ResponseWriter, r *request) error {

	photo, err := c.ds.photos.get(r.getIntParam("id"))
	if err != nil {
		return err
	}

	if !photo.canDelete(r.user) {
		return httpError{http.StatusForbidden, "You're not allowed to delete this photo"}
	}
	if err := c.ds.photos.remove(photo); err != nil {
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

	sendMessage(&socketMessage{r.user.Name, "", photo.ID, "photo_deleted"})
	return renderString(w, http.StatusOK, "Photo deleted")
}

func getPhotoDetail(c *appContext, w http.ResponseWriter, r *request) error {

	user, err := c.authenticate(r, false)
	if err != nil {
		return err
	}

	photo, err := c.ds.photos.getDetail(r.getIntParam("id"), user)
	if err != nil {
		return err
	}
	return renderJSON(w, photo, http.StatusOK)

}

func getPhotoToEdit(c *appContext, w http.ResponseWriter, r *request) (*photo, error) {
	user, err := c.authenticate(r, true)
	if err != nil {
		return nil, err
	}

	photo, err := c.ds.photos.get(r.getIntParam("id"))
	if err != nil {
		return photo, err
	}

	if !photo.canEdit(user) {
		return photo, httpError{http.StatusForbidden, "You're not allowed to edit this photo"}
	}
	return photo, nil
}

func editPhotoTitle(c *appContext, w http.ResponseWriter, r *request) error {

	photo, err := getPhotoToEdit(c, w, r)

	if err != nil {
		return err
	}

	s := &struct {
		Title string `json:"title"`
	}{}

	if err := r.decodeJSON(s); err != nil {
		return err
	}

	photo.Title = s.Title

	if err := validate(newPhotoValidator(photo)); err != nil {
		return err

	}

	if err := c.ds.photos.update(photo); err != nil {
		return err
	}
	sendMessage(&socketMessage{r.user.Name, "", photo.ID, "photo_updated"})
	return renderString(w, http.StatusOK, "Photo updated")
}

func editPhotoTags(c *appContext, w http.ResponseWriter, r *request) error {

	photo, err := getPhotoToEdit(c, w, r)
	if err != nil {
		return err
	}

	s := &struct {
		Tags []string `json:"tags"`
	}{}

	if err := r.decodeJSON(s); err != nil {
		return err
	}

	photo.Tags = s.Tags

	if err := c.ds.photos.updateTags(photo); err != nil {
		return err
	}
	sendMessage(&socketMessage{r.user.Name, "", photo.ID, "photo_updated"})
	return renderString(w, http.StatusOK, "Photo updated")

}

func upload(c *appContext, w http.ResponseWriter, r *request) error {

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
		OwnerID:  r.user.ID,
		Filename: filename,
		Tags:     tags,
	}

	if err := validate(newPhotoValidator(photo)); err != nil {
		return err
	}

	if err := c.ds.photos.create(photo); err != nil {
		return err
	}

	if err := c.cache.clear(); err != nil {
		return err
	}

	sendMessage(&socketMessage{r.user.Name, "", photo.ID, "photo_uploaded"})
	return renderJSON(w, photo, http.StatusCreated)
}

func searchPhotos(c *appContext, w http.ResponseWriter, r *request) error {

	page := r.getPage()
	q := r.FormValue("q")
	qKey := base64.StdEncoding.EncodeToString([]byte(q))
	cacheKey := fmt.Sprintf("photos:search:%s:page:%d", qKey, page.index)

	return c.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := c.ds.photos.search(page, q)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})

}

func photosByOwnerID(c *appContext, w http.ResponseWriter, r *request) error {

	page := r.getPage()
	ownerID := r.getIntParam("ownerID")
	cacheKey := fmt.Sprintf("photos:ownerID:%d:page:%d", ownerID, page.index)

	return c.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := c.ds.photos.byOwnerID(page, ownerID)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func getPhotos(c *appContext, w http.ResponseWriter, r *request) error {

	page := r.getPage()
	orderBy := r.FormValue("orderBy")
	cacheKey := fmt.Sprintf("photos:%s:page:%d", orderBy, page.index)

	return c.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := c.ds.photos.all(page, orderBy)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func getTags(c *appContext, w http.ResponseWriter, r *request) error {
	return c.cache.render(w, http.StatusOK, "tags", func() (interface{}, error) {
		tags, err := c.ds.photos.getTagCounts()
		if err != nil {
			return tags, err
		}
		return tags, nil
	})

}

func voteDown(c *appContext, w http.ResponseWriter, r *request) error {
	return vote(c, w, r, func(photo *photo) { photo.DownVotes++ })
}

func voteUp(c *appContext, w http.ResponseWriter, r *request) error {
	return vote(c, w, r, func(photo *photo) { photo.UpVotes++ })
}

func vote(c *appContext, w http.ResponseWriter, r *request, fn func(photo *photo)) error {
	var (
		photo *photo
		err   error
	)
	photo, err = c.ds.photos.get(r.getIntParam("id"))
	if err != nil {
		return err
	}

	if !photo.canVote(r.user) {
		return httpError{http.StatusForbidden, "You're not allowed to vote on this photo"}
	}

	fn(photo)

	if err = c.ds.photos.update(photo); err != nil {
		return err
	}

	r.user.registerVote(photo.ID)

	if err = c.ds.users.update(r.user); err != nil {
		return err
	}
	return renderString(w, http.StatusOK, "Voting successful")
}
