package api

import (
	"encoding/base64"
	"fmt"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"strings"
)

func (a *appContext) deletePhoto(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, true)
	if err != nil {
		return err
	}
	photo, err := a.ds.photos.get(getIntParam(c, "id"))
	if err != nil {
		return err
	}

	if !photo.canDelete(user) {
		return httpError(http.StatusForbidden, "You're not allowed to delete this photo")

	}
	if err := a.ds.photos.delete(photo); err != nil {
		return err
	}

	go func() {
		if err := a.fs.clean(photo.Filename); err != nil {
			log.Println(err)
		}
	}()

	if err := a.cache.clear(); err != nil {
		return err
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_deleted"})
	return renderString(w, http.StatusOK, "Photo deleted")
}

func (a *appContext) photoDetail(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, false)
	if err != nil {
		return err
	}

	photo, err := a.ds.photos.getDetail(getIntParam(c, "id"), user)
	if err != nil {
		return err
	}
	return renderJSON(w, photo, http.StatusOK)

}

func (a *appContext) getPhotoToEdit(c web.C, w http.ResponseWriter, r *http.Request) (*photo, error) {
	user, err := a.authenticate(c, r, true)
	if err != nil {
		return nil, err
	}

	photo, err := a.ds.photos.get(getIntParam(c, "id"))
	if err != nil {
		return photo, err
	}

	if !photo.canEdit(user) {
		return photo, httpError(http.StatusForbidden, "You're not allowed to edit this photo")
	}
	return photo, nil
}

func (a *appContext) editPhotoTitle(c web.C, w http.ResponseWriter, r *http.Request) error {

	photo, err := a.getPhotoToEdit(c, w, r)

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

	if err := validate(newPhotoValidator(photo)); err != nil {
		return err

	}

	if err := a.ds.photos.update(photo); err != nil {
		return err
	}
	if user, err := a.authenticate(c, r, true); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	return renderString(w, http.StatusOK, "Photo updated")
}

func (a *appContext) editPhotoTags(c web.C, w http.ResponseWriter, r *http.Request) error {

	photo, err := a.getPhotoToEdit(c, w, r)
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

	if err := a.ds.photos.updateTags(photo); err != nil {
		return err
	}
	if user, err := a.authenticate(c, r, true); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	return renderString(w, http.StatusOK, "Photo updated")

}

func (a *appContext) upload(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, true)
	if err != nil {
		return err
	}

	title := r.FormValue("title")
	taglist := r.FormValue("taglist")
	tags := strings.Split(taglist, " ")

	src, hdr, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			return httpError(http.StatusBadRequest, "Invalid photo")
		}
		return err
	}
	defer src.Close()

	contentType := hdr.Header["Content-Type"][0]

	filename, err := a.fs.store(src, contentType)

	if err != nil {
		if err == errInvalidContentType {
			return httpError(http.StatusBadRequest, err.Error())
		}
		return err
	}

	photo := &photo{Title: title,
		OwnerID:  user.ID,
		Filename: filename,
		Tags:     tags,
	}

	if err := validate(newPhotoValidator(photo)); err != nil {
		return err
	}

	if err := a.ds.photos.insert(photo); err != nil {
		return err
	}

	if err := a.cache.clear(); err != nil {
		return err
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_uploaded"})
	return renderJSON(w, photo, http.StatusCreated)
}

func (a *appContext) searchPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	q := r.FormValue("q")
	qKey := base64.StdEncoding.EncodeToString([]byte(q))
	cacheKey := fmt.Sprintf("photos:search:%s:page:%d", qKey, page.index)

	return a.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := a.ds.photos.search(page, q)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})

}

func (a *appContext) photosByOwnerID(c web.C, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	ownerID := getIntParam(c, "ownerID")
	cacheKey := fmt.Sprintf("photos:ownerID:%d:page:%d", ownerID, page.index)

	return a.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := a.ds.photos.byOwnerID(page, ownerID)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func (a *appContext) getPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	orderBy := r.FormValue("orderBy")
	cacheKey := fmt.Sprintf("photos:%s:page:%d", orderBy, page.index)

	return a.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := a.ds.photos.all(page, orderBy)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func (a *appContext) getTags(_ web.C, w http.ResponseWriter, r *http.Request) error {
	return a.cache.render(w, http.StatusOK, "tags", func() (interface{}, error) {
		tags, err := a.ds.photos.getTagCounts()
		if err != nil {
			return tags, err
		}
		return tags, nil
	})

}

func (a *appContext) voteDown(c web.C, w http.ResponseWriter, r *http.Request) error {
	return a.vote(c, w, r, func(photo *photo) { photo.DownVotes++ })
}

func (a *appContext) voteUp(c web.C, w http.ResponseWriter, r *http.Request) error {
	return a.vote(c, w, r, func(photo *photo) { photo.UpVotes++ })
}

func (a *appContext) vote(c web.C, w http.ResponseWriter, r *http.Request, fn func(photo *photo)) error {
	var (
		photo *photo
		err   error
	)
	user, err := a.authenticate(c, r, true)

	if err != nil {
		return err
	}

	photo, err = a.ds.photos.get(getIntParam(c, "id"))
	if err != nil {
		return err
	}

	if !photo.canVote(user) {
		return httpError(http.StatusForbidden, "You're not allowed to vote on this photo")
	}

	fn(photo)

	if err = a.ds.photos.update(photo); err != nil {
		return err
	}

	user.registerVote(photo.ID)

	if err = a.ds.users.update(user); err != nil {
		return err
	}
	return renderString(w, http.StatusOK, "Voting successful")
}
