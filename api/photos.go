package api

import (
	"encoding/base64"
	"fmt"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"strings"
)

func (a *AppContext) deletePhoto(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, true)
	if err != nil {
		return err
	}
	photo, err := a.photoDS.Get(getIntParam(c, "id"))
	if err != nil {
		return err
	}

	if !photo.CanDelete(user) {
		return httpError(http.StatusForbidden, "You're not allowed to delete this photo")

	}
	if err := a.photoDS.Delete(photo); err != nil {
		return err
	}

	go func() {
		if err := a.fs.Clean(photo.Filename); err != nil {
			log.Println(err)
		}
	}()

	if err := a.cache.DeleteAll(); err != nil {
		return err
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_deleted"})
	return renderString(w, http.StatusOK, "Photo deleted")
}

func (a *AppContext) photoDetail(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, false)
	if err != nil {
		return err
	}

	photo, err := a.photoDS.GetDetail(getIntParam(c, "id"), user)
	if err != nil {
		return err
	}
	return renderJSON(w, photo, http.StatusOK)

}

func (a *AppContext) getPhotoToEdit(c web.C, w http.ResponseWriter, r *http.Request) (*Photo, error) {
	user, err := a.authenticate(c, r, true)
	if err != nil {
		return nil, err
	}

	photo, err := a.photoDS.Get(getIntParam(c, "id"))
	if err != nil {
		return photo, err
	}

	if !photo.CanEdit(user) {
		return photo, httpError(http.StatusForbidden, "You're not allowed to edit this photo")
	}
	return photo, nil
}

func (a *AppContext) editPhotoTitle(c web.C, w http.ResponseWriter, r *http.Request) error {

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

	if err := validate(NewPhotoValidator(photo)); err != nil {
		return err
	}

	if err := a.photoDS.Update(photo); err != nil {
		return err
	}
	if user, err := a.authenticate(c, r, true); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	return renderString(w, http.StatusOK, "Photo updated")
}

func (a *AppContext) editPhotoTags(c web.C, w http.ResponseWriter, r *http.Request) error {

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

	if err := a.photoDS.UpdateTags(photo); err != nil {
		return err
	}
	if user, err := a.authenticate(c, r, true); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	return renderString(w, http.StatusOK, "Photo updated")

}

func (a *AppContext) upload(c web.C, w http.ResponseWriter, r *http.Request) error {

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

	filename, err := a.fs.Store(src, contentType)

	if err != nil {
		if err == InvalidContentType {
			return httpError(http.StatusBadRequest, err.Error())
		}
		return err
	}

	photo := &Photo{Title: title,
		OwnerID:  user.ID,
		Filename: filename,
		Tags:     tags,
	}

	if err := validate(NewPhotoValidator(photo)); err != nil {
		return err
	}

	if err := a.photoDS.Insert(photo); err != nil {
		return err
	}

	if err := a.cache.DeleteAll(); err != nil {
		return err
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_uploaded"})
	return renderJSON(w, photo, http.StatusCreated)
}

func (a *AppContext) searchPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	q := r.FormValue("q")
	qKey := base64.StdEncoding.EncodeToString([]byte(q))
	cacheKey := fmt.Sprintf("photos:search:%s:page:%d", qKey, page.Index)

	return a.cache.Render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := a.photoDS.Search(page, q)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})

}

func (a *AppContext) photosByOwnerID(c web.C, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	ownerID := getIntParam(c, "ownerID")
	cacheKey := fmt.Sprintf("photos:ownerID:%d:page:%d", ownerID, page.Index)

	return a.cache.Render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := a.photoDS.ByOwnerID(page, ownerID)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func (a *AppContext) getPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	orderBy := r.FormValue("orderBy")
	cacheKey := fmt.Sprintf("photos:%s:page:%d", orderBy, page.Index)

	return a.cache.Render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := a.photoDS.All(page, orderBy)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func (a *AppContext) getTags(_ web.C, w http.ResponseWriter, r *http.Request) error {
	return a.cache.Render(w, http.StatusOK, "tags", func() (interface{}, error) {
		tags, err := a.photoDS.GetTagCounts()
		if err != nil {
			return tags, err
		}
		return tags, nil
	})

}

func (a *AppContext) voteDown(c web.C, w http.ResponseWriter, r *http.Request) error {
	return a.vote(c, w, r, func(photo *Photo) { photo.DownVotes += 1 })
}

func (a *AppContext) voteUp(c web.C, w http.ResponseWriter, r *http.Request) error {
	return a.vote(c, w, r, func(photo *Photo) { photo.UpVotes += 1 })
}

func (a *AppContext) vote(c web.C, w http.ResponseWriter, r *http.Request, fn func(photo *Photo)) error {
	var (
		photo *Photo
		err   error
	)
	user, err := a.authenticate(c, r, true)

	if err != nil {
		return err
	}

	photo, err = a.photoDS.Get(getIntParam(c, "id"))
	if err != nil {
		return err
	}

	if !photo.CanVote(user) {
		return httpError(http.StatusForbidden, "You're not allowed to vote on this photo")
	}

	fn(photo)

	if err = a.photoDS.Update(photo); err != nil {
		return err
	}

	user.RegisterVote(photo.ID)

	if err = a.userDS.Update(user); err != nil {
		return err
	}
	return renderString(w, http.StatusOK, "Voting successful")
}
