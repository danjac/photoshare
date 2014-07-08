package api

import (
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
	"strings"
)

var (
	allowedContentTypes = []string{"image/png", "image/jpeg"}
)

func isAllowedContentType(contentType string) bool {
	for _, value := range allowedContentTypes {
		if contentType == value {
			return true
		}
	}

	return false
}

func getPage(r *http.Request) int64 {
	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		page = 1
	}
	return page
}

func getPhotoDetail(c web.C, user *User) (*PhotoDetail, error) {
	photoID, err := strconv.ParseInt(c.URLParams["id"], 10, 0)
	if err != nil {
		return nil, nil
	}
	return photoMgr.GetDetail(photoID, user)
}

func getPhoto(c web.C) (*Photo, error) {
	photoID, err := strconv.ParseInt(c.URLParams["id"], 10, 0)
	if err != nil {
		return nil, nil
	}
	return photoMgr.Get(photoID)
}

func deletePhoto(c web.C, w http.ResponseWriter, r *http.Request) {

	user, ok := checkAuth(c, w, r)
	if !ok {
		return
	}

	photo, err := getPhoto(c)
	if err != nil {
		handleServerError(w, err)
		return
	}

	if photo == nil {
		http.NotFound(w, r)
		return
	}
	if !photo.CanDelete(user) {
		http.Error(w, "You can't delete this photo", http.StatusForbidden)
		return
	}
	if err := photoMgr.Delete(photo); err != nil {
		handleServerError(w, err)
		return
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_deleted"})
	w.WriteHeader(http.StatusOK)
}

func photoDetail(c web.C, w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(c, r)
	if err != nil {
		handleServerError(w, err)
		return
	}

	photo, err := getPhotoDetail(c, user)
	if err != nil {
		handleServerError(w, err)
		return
	}
	if photo == nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, photo, http.StatusOK)
}

func getPhotoToEdit(c web.C, w http.ResponseWriter, r *http.Request) (*Photo, bool) {
	user, ok := checkAuth(c, w, r)
	if !ok {
		return nil, false
	}

	photo, err := getPhoto(c)

	if err != nil {
		handleServerError(w, err)
		return nil, false
	}

	if photo == nil {
		http.NotFound(w, r)
		return nil, false
	}

	if !photo.CanEdit(user) {
		http.Error(w, "You can't edit this photo", http.StatusForbidden)
		return photo, false
	}
	return photo, true
}

func editPhotoTitle(c web.C, w http.ResponseWriter, r *http.Request) {

	photo, ok := getPhotoToEdit(c, w, r)

	if !ok {
		return
	}

	s := &struct {
		Title string `json:"title"`
	}{}

	if err := parseJSON(r, s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	photo.Title = s.Title

	validator := getPhotoValidator(photo)

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			handleServerError(w, err)
			return
		}
		writeJSON(w, result, http.StatusBadRequest)
		return
	}

	if err := photoMgr.Update(photo); err != nil {
		handleServerError(w, err)
		return
	}
	if user, err := getCurrentUser(c, r); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	w.WriteHeader(http.StatusOK)
}

func editPhotoTags(c web.C, w http.ResponseWriter, r *http.Request) {

	photo, ok := getPhotoToEdit(c, w, r)

	if !ok {
		return
	}

	s := &struct {
		Tags []string `json:"tags"`
	}{}

	if err := parseJSON(r, s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	photo.Tags = s.Tags

	if err := photoMgr.UpdateTags(photo); err != nil {
		handleServerError(w, err)
		return
	}
	if user, err := getCurrentUser(c, r); err == nil {
		sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_updated"})
	}
	w.WriteHeader(http.StatusOK)
}

func upload(c web.C, w http.ResponseWriter, r *http.Request) {

	user, ok := checkAuth(c, w, r)
	if !ok {
		return
	}

	title := r.FormValue("title")
	taglist := r.FormValue("taglist")
	tags := strings.Split(taglist, " ")

	src, hdr, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			http.Error(w, "No image was posted", http.StatusBadRequest)
			return
		}
		handleServerError(w, err)
		return
	}
	contentType := hdr.Header["Content-Type"][0]

	if !isAllowedContentType(contentType) {
		http.Error(w, "No image was posted", http.StatusBadRequest)
		return
	}

	defer src.Close()

	filename, err := imageProcessor.Process(src, contentType)

	if err != nil {
		handleServerError(w, err)
		return
	}

	photo := &Photo{Title: title,
		OwnerID:  user.ID,
		Filename: filename,
		Tags:     tags,
	}

	validator := getPhotoValidator(photo)

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			handleServerError(w, err)
			return
		}
		writeJSON(w, result, http.StatusBadRequest)
		return
	}

	if err := photoMgr.Insert(photo); err != nil {
		handleServerError(w, err)
		return
	}

	sendMessage(&SocketMessage{user.Name, "", photo.ID, "photo_uploaded"})
	writeJSON(w, photo, http.StatusOK)
}

func searchPhotos(c web.C, w http.ResponseWriter, r *http.Request) {
	photos, err := photoMgr.Search(getPage(r), r.FormValue("q"))
	if err != nil {
		handleServerError(w, err)
		return
	}
	writeJSON(w, photos, http.StatusOK)
}

func photosByOwnerID(c web.C, w http.ResponseWriter, r *http.Request) {
	ownerID, err := strconv.ParseInt(c.URLParams["ownerID"], 10, 0)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	photos, err := photoMgr.ByOwnerID(getPage(r), ownerID)
	if err != nil {
		handleServerError(w, err)
		return
	}
	writeJSON(w, photos, http.StatusOK)
}

func getPhotos(c web.C, w http.ResponseWriter, r *http.Request) {
	photos, err := photoMgr.All(getPage(r), r.FormValue("orderBy"))
	if err != nil {
		handleServerError(w, err)
		return
	}
	writeJSON(w, photos, http.StatusOK)
}

func getTags(c web.C, w http.ResponseWriter, r *http.Request) {
	tags, err := photoMgr.GetTagCounts()
	if err != nil {
		handleServerError(w, err)
		return
	}
	writeJSON(w, tags, http.StatusOK)
}

func voteDown(c web.C, w http.ResponseWriter, r *http.Request) {
	vote(c, w, r, func(photo *Photo) { photo.DownVotes += 1 })
}

func voteUp(c web.C, w http.ResponseWriter, r *http.Request) {
	vote(c, w, r, func(photo *Photo) { photo.UpVotes += 1 })
}

func vote(c web.C, w http.ResponseWriter, r *http.Request, fn func(photo *Photo)) {
	var (
		photo *Photo
		err   error
	)
	user, ok := checkAuth(c, w, r)
	if !ok {
		return
	}

	photo, err = getPhoto(c)
	if err != nil {
		handleServerError(w, err)
		return
	}
	if photo == nil {
		http.NotFound(w, r)
		return
	}

	if !photo.CanVote(user) {
		http.Error(w, "You can't vote on this photo", http.StatusForbidden)
		return
	}

	fn(photo)

	if err = photoMgr.Update(photo); err != nil {
		handleServerError(w, err)
		return
	}

	user.RegisterVote(photo.ID)

	if err = userMgr.Update(user); err != nil {
		handleServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
