package api

import (
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (a *AppContext) deletePhoto(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, true)
	if err != nil {
		return err
	}
	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err := a.ds.GetPhoto(photoID)
	if err != nil {
		return err
	}

	if !photo.CanDelete(user) {
		return httpError(http.StatusForbidden, "You're not allowed to delete this photo")
	}
	tx, err := a.ds.Begin()
	if err != nil {
		return err
	}
	if err := tx.DeletePhoto(photo); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		if err := a.fileMgr.Clean(photo.Filename); err != nil {
			log.Println(err)
		}
	}()

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_deleted"})
	return renderString(w, http.StatusOK, "Photo deleted")
}

func (a *AppContext) photoDetail(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, false)
	if err != nil {
		return err
	}

	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err := a.ds.GetPhotoDetail(photoID, user)
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

	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err := a.ds.GetPhoto(photoID)
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

	validator := NewPhotoValidator(photo)

	if err := validate(validator); err != nil {
		return err
	}

	tx, err := a.ds.Begin()
	if err != nil {
		return err
	}
	if err := tx.UpdatePhoto(photo); err != nil {
		return err
	}
	tx.Commit()

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

	tx, err := a.ds.Begin()
	if err != nil {
		return err
	}
	if err := tx.UpdateTags(photo); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
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

	filename, err := a.fileMgr.Store(src, contentType)

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

	validator := NewPhotoValidator(photo)

	if err := validate(validator); err != nil {
		return err
	}

	tx, err := a.ds.Begin()
	if err != nil {
		return err
	}

	if err := tx.InsertPhoto(photo); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.UpdateTags(photo); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_uploaded"})
	return renderJSON(w, photo, http.StatusCreated)
}

func (a *AppContext) searchPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {
	photos, err := a.ds.SearchPhotos(getPage(r), r.FormValue("q"))
	if err != nil {
		return err
	}
	return renderJSON(w, photos, http.StatusOK)
}

func (a *AppContext) photosByOwnerID(c web.C, w http.ResponseWriter, r *http.Request) error {
	ownerID, err := strconv.ParseInt(c.URLParams["ownerID"], 10, 0)
	if err != nil {
		return err
	}
	photos, err := a.ds.GetPhotosByOwnerID(getPage(r), ownerID)
	if err != nil {
		return err
	}
	return renderJSON(w, photos, http.StatusOK)
}

func (a *AppContext) getPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {
	photos, err := a.ds.GetPhotos(getPage(r), r.FormValue("orderBy"))
	if err != nil {
		return err
	}
	return renderJSON(w, photos, http.StatusOK)
}

func (a *AppContext) getTags(_ web.C, w http.ResponseWriter, r *http.Request) error {
	tags, err := a.ds.GetTagCounts()
	if err != nil {
		return err
	}
	return renderJSON(w, tags, http.StatusOK)
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

	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err = a.ds.GetPhoto(photoID)
	if err != nil {
		return err
	}

	if !photo.CanVote(user) {
		return httpError(http.StatusForbidden, "You're not allowed to vote on this photo")
	}

	fn(photo)

	tx, err := a.ds.Begin()
	if err != nil {
		return err
	}

	if err = tx.UpdatePhoto(photo); err != nil {
		tx.Rollback()
		return err
	}

	user.RegisterVote(photo.ID)

	if err = tx.UpdateUser(user); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return renderString(w, http.StatusOK, "Voting successful")
}
