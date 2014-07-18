package api

import (
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
	"strings"
)

func deletePhoto(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := getCurrentUser(r, true)
	if err != nil {
		return err
	}
	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err := photoMgr.Get(photoID)
	if err != nil {
		return err
	}

	if !photo.CanDelete(user) {
		return HttpError{http.StatusForbidden}
	}
	if err := photoMgr.Delete(photo); err != nil {
		return err
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_deleted"})
	return renderStatus(w, http.StatusNoContent)
}

func photoDetail(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := getCurrentUser(r, false)
	if err != nil {
		return err
	}

	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err := photoMgr.GetDetail(photoID, user)
	if err != nil {
		return err
	}
	return renderJSON(w, photo, http.StatusOK)
}

func getPhotoToEdit(c web.C, w http.ResponseWriter, r *http.Request) (*Photo, error) {
	user, err := getCurrentUser(r, true)
	if err != nil {
		return nil, err
	}

	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err := photoMgr.Get(photoID)
	if err != nil {
		return photo, err
	}

	if !photo.CanEdit(user) {
		return photo, HttpError{http.StatusForbidden}
	}
	return photo, nil
}

func editPhotoTitle(c web.C, w http.ResponseWriter, r *http.Request) error {

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

	validator := getPhotoValidator(photo)

	if err := formHandler.Validate(validator); err != nil {
		return err
	}

	if err := photoMgr.Update(photo); err != nil {
		return err
	}
	if user, err := getCurrentUser(r, true); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	return renderStatus(w, http.StatusNoContent)
}

func editPhotoTags(c web.C, w http.ResponseWriter, r *http.Request) error {

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

	if err := photoMgr.UpdateTags(photo); err != nil {
		return err
	}
	if user, err := getCurrentUser(r, true); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	w.WriteHeader(http.StatusNoContent)
	return nil

}

func upload(_ web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := getCurrentUser(r, true)
	if err != nil {
		return err
	}

	title := r.FormValue("title")
	taglist := r.FormValue("taglist")
	tags := strings.Split(taglist, " ")

	src, hdr, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			return HttpError{http.StatusBadRequest}
		}
		return err
	}
	defer src.Close()

	contentType := hdr.Header["Content-Type"][0]

	filename, err := imageProcessor.Process(src, contentType)

	if err != nil {
		if err == InvalidContentType {
			return HttpError{http.StatusBadRequest}
		}
		return err
	}

	photo := &Photo{Title: title,
		OwnerID:  user.ID,
		Filename: filename,
		Tags:     tags,
	}

	validator := getPhotoValidator(photo)

	if err := formHandler.Validate(validator); err != nil {
		return err
	}

	if err := photoMgr.Insert(photo); err != nil {
		return err
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_uploaded"})
	return renderJSON(w, photo, http.StatusCreated)
}

func searchPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {
	photos, err := photoMgr.Search(getPage(r), r.FormValue("q"))
	if err != nil {
		return err
	}
	return renderJSON(w, photos, http.StatusOK)
}

func photosByOwnerID(c web.C, w http.ResponseWriter, r *http.Request) error {
	ownerID, err := strconv.ParseInt(c.URLParams["ownerID"], 10, 0)
	if err != nil {
		return err
	}
	photos, err := photoMgr.ByOwnerID(getPage(r), ownerID)
	if err != nil {
		return err
	}
	return renderJSON(w, photos, http.StatusOK)
}

func getPhotos(_ web.C, w http.ResponseWriter, r *http.Request) error {
	photos, err := photoMgr.All(getPage(r), r.FormValue("orderBy"))
	if err != nil {
		return err
	}
	return renderJSON(w, photos, http.StatusOK)
}

func getTags(_ web.C, w http.ResponseWriter, r *http.Request) error {
	tags, err := photoMgr.GetTagCounts()
	if err != nil {
		return err
	}
	return renderJSON(w, tags, http.StatusOK)
}

func voteDown(c web.C, w http.ResponseWriter, r *http.Request) error {
	return vote(c, w, r, func(photo *Photo) { photo.DownVotes += 1 })
}

func voteUp(c web.C, w http.ResponseWriter, r *http.Request) error {
	return vote(c, w, r, func(photo *Photo) { photo.UpVotes += 1 })
}

func vote(c web.C, w http.ResponseWriter, r *http.Request, fn func(photo *Photo)) error {
	var (
		photo *Photo
		err   error
	)
	user, err := getCurrentUser(r, true)

	if err != nil {
		return err
	}

	photoID, _ := strconv.ParseInt(c.URLParams["id"], 10, 0)
	photo, err = photoMgr.Get(photoID)
	if err != nil {
		return err
	}

	if !photo.CanVote(user) {
		return HttpError{http.StatusForbidden}
	}

	fn(photo)

	if err = photoMgr.Update(photo); err != nil {
		return err
	}

	user.RegisterVote(photo.ID)

	if err = userMgr.Update(user); err != nil {
		return err
	}
	return renderStatus(w, http.StatusNoContent)
}
