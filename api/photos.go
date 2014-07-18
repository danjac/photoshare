package api

import (
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (a *AppContext) deletePhoto(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(r, true)
	if err != nil {
		return err
	}
	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err := a.photoMgr.Get(photoID)
	if err != nil {
		return err
	}

	if !photo.CanDelete(user) {
		return httpError(http.StatusForbidden, "You're not allowed to delete this photo")
	}
	if err := a.photoMgr.Delete(photo); err != nil {
		return err
	}

	go func() {
		if err := a.fileMgr.Clean(photo.Filename); err != nil {
			log.Println(err)
		}
	}()

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_deleted"})
	return renderStatus(w, http.StatusOK, "Photo deleted")
}

func (a *AppContext) photoDetail(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(r, false)
	if err != nil {
		return err
	}

	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err := a.photoMgr.GetDetail(photoID, user)
	if err != nil {
		return err
	}
	return renderJSON(w, photo, http.StatusOK)
}

func (a *AppContext) getPhotoToEdit(c web.C, w http.ResponseWriter, r *http.Request) (*Photo, error) {
	user, err := a.authenticate(r, true)
	if err != nil {
		return nil, err
	}

	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err := a.photoMgr.Get(photoID)
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

	validator := getPhotoValidator(photo)

	if err := validate(validator); err != nil {
		return err
	}

	if err := a.photoMgr.Update(photo); err != nil {
		return err
	}
	if user, err := a.authenticate(r, true); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	return renderStatus(w, http.StatusOK, "Photo updated")
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

	if err := a.photoMgr.UpdateTags(photo); err != nil {
		return err
	}
	if user, err := a.authenticate(r, true); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	return renderStatus(w, http.StatusOK, "Photo updated")

}

func (a *AppContext) upload(_ web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(r, true)
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

	validator := getPhotoValidator(photo)

	if err := validate(validator); err != nil {
		return err
	}

	if err := a.photoMgr.Insert(photo); err != nil {
		return err
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_uploaded"})
	return renderJSON(w, photo, http.StatusCreated)
}

func (a *AppContext) searchPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {
	photos, err := a.photoMgr.Search(getPage(r), r.FormValue("q"))
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
	photos, err := a.photoMgr.ByOwnerID(getPage(r), ownerID)
	if err != nil {
		return err
	}
	return renderJSON(w, photos, http.StatusOK)
}

func (a *AppContext) getPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {
	photos, err := a.photoMgr.All(getPage(r), r.FormValue("orderBy"))
	if err != nil {
		return err
	}
	return renderJSON(w, photos, http.StatusOK)
}

func (a *AppContext) getTags(_ web.C, w http.ResponseWriter, r *http.Request) error {
	tags, err := a.photoMgr.GetTagCounts()
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
	user, err := a.authenticate(r, true)

	if err != nil {
		return err
	}

	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err = a.photoMgr.Get(photoID)
	if err != nil {
		return err
	}

	if !photo.CanVote(user) {
		return httpError(http.StatusForbidden, "You're not allowed to vote on this photo")
	}

	fn(photo)

	if err = a.photoMgr.Update(photo); err != nil {
		return err
	}

	user.RegisterVote(photo.ID)

	if err = a.userMgr.Update(user); err != nil {
		return err
	}
	return renderStatus(w, http.StatusOK, "Voting successful")
}
