package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/storage"
	"github.com/danjac/photoshare/api/validation"
	"github.com/zenazn/goji/web"
	"net/http"
	"strconv"
	"strings"
)

var (
	allowedContentTypes = []string{"image/png", "image/jpeg"}
	imageProcessor      = storage.NewImageProcessor()
)

var getPhotoValidator = func(photo *models.Photo) validation.Validator {
	return validation.NewPhotoValidator(photo)
}

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

func getPhotoDetail(c web.C, user *models.User) (*models.PhotoDetail, error) {
	photoID, err := strconv.ParseInt(c.URLParams["id"], 10, 0)
	if err != nil {
		return nil, nil
	}
	return photoMgr.GetDetail(photoID, user)
}

func getPhoto(c web.C) (*models.Photo, error) {
	photoID, err := strconv.ParseInt(c.URLParams["id"], 10, 0)
	if err != nil {
		return nil, nil
	}
	return photoMgr.Get(photoID)
}

func deletePhoto(c web.C, w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(c, r)

	if err != nil {
		panic(err)
	}

	if !user.IsAuthenticated {
		writeError(w, http.StatusUnauthorized)
		return
	}

	photo, err := getPhoto(c)
	if err != nil {
		panic(err)
	}

	if photo == nil {
		http.NotFound(w, r)
		return
	}
	if !photo.CanDelete(user) {
		writeError(w, http.StatusForbidden)
		return
	}
	if err := photoMgr.Delete(photo); err != nil {
		panic(err)
	}

	sendMessage(&Message{user.Name, "", photo.ID, "photo_deleted"})
	w.WriteHeader(http.StatusOK)
}

func photoDetail(c web.C, w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(c, r)
	if err != nil {
		panic(err)
	}

	photo, err := getPhotoDetail(c, user)
	if err != nil {
		panic(err)
	}
	if photo == nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, photo, http.StatusOK)
}

func getPhotoToEdit(c web.C, w http.ResponseWriter, r *http.Request) (*models.Photo, bool) {
	user, err := getCurrentUser(c, r)
	if err != nil {
		panic(err)
	}

	if !user.IsAuthenticated {
		writeError(w, http.StatusUnauthorized)
		return nil, false
	}

	photo, err := getPhoto(c)

	if err != nil {
		panic(err)
	}

	if photo == nil {
		http.NotFound(w, r)
		return nil, false
	}

	if !photo.CanEdit(user) {
		writeError(w, http.StatusForbidden)
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
			panic(err)
		}
		writeJSON(w, result, http.StatusBadRequest)
		return
	}

	if err := photoMgr.Update(photo); err != nil {
		panic(err)
	}
	if user, err := getCurrentUser(c, r); err == nil {
		sendMessage(&Message{user.Name, "", photo.ID, "photo_updated"})
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
		panic(err)
	}
	if user, err := getCurrentUser(c, r); err == nil {
		sendMessage(&Message{user.Name, "", photo.ID, "photo_updated"})
	}
	w.WriteHeader(http.StatusOK)
}

func upload(c web.C, w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(c, r)
	if err != nil {
		panic(err)
	}
	if !user.IsAuthenticated {
		writeError(w, http.StatusUnauthorized)
		return
	}
	title := r.FormValue("title")
	taglist := r.FormValue("taglist")
	tags := strings.Split(taglist, " ")

	src, hdr, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			writeString(w, "No image was posted", http.StatusBadRequest)
			return
		}
		panic(err)
	}
	contentType := hdr.Header["Content-Type"][0]

	if !isAllowedContentType(contentType) {
		writeString(w, "No image was posted", http.StatusBadRequest)
		return
	}

	defer src.Close()

	filename, err := imageProcessor.Process(src, contentType)

	if err != nil {
		panic(err)
	}

	photo := &models.Photo{Title: title,
		OwnerID:  user.ID,
		Filename: filename,
		Tags:     tags,
	}

	validator := getPhotoValidator(photo)

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			panic(err)
		}
		writeJSON(w, result, http.StatusBadRequest)
		return
	}

	if err := photoMgr.Insert(photo); err != nil {
		panic(err)
	}

	sendMessage(&Message{user.Name, "", photo.ID, "photo_uploaded"})
	writeJSON(w, photo, http.StatusOK)
}

func searchPhotos(c web.C, w http.ResponseWriter, r *http.Request) {
	photos, err := photoMgr.Search(getPage(r), r.FormValue("q"))
	if err != nil {
		panic(err)
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
		panic(err)
	}
	writeJSON(w, photos, http.StatusOK)
}

func getPhotos(c web.C, w http.ResponseWriter, r *http.Request) {
	photos, err := photoMgr.All(getPage(r), r.FormValue("orderBy"))
	if err != nil {
		panic(err)
	}
	writeJSON(w, photos, http.StatusOK)
}

func getTags(c web.C, w http.ResponseWriter, r *http.Request) {
	tags, err := photoMgr.GetTagCounts()
	if err != nil {
		panic(err)
	}
	writeJSON(w, tags, http.StatusOK)
}

func voteDown(c web.C, w http.ResponseWriter, r *http.Request) {
	vote(c, w, r, func(photo *models.Photo) { photo.DownVotes += 1 })
}

func voteUp(c web.C, w http.ResponseWriter, r *http.Request) {
	vote(c, w, r, func(photo *models.Photo) { photo.UpVotes += 1 })
}

func vote(c web.C, w http.ResponseWriter, r *http.Request, fn func(photo *models.Photo)) {
	var (
		photo *models.Photo
		err   error
	)

	user, err := getCurrentUser(c, r)
	if err != nil {
		panic(err)
	}

	if !user.IsAuthenticated {
		writeError(w, http.StatusUnauthorized)
		return
	}
	photo, err = getPhoto(c)
	if err != nil {
		panic(err)
	}
	if photo == nil {
		http.NotFound(w, r)
		return
	}

	if !photo.CanVote(user) {
		writeError(w, http.StatusForbidden)
		return
	}

	fn(photo)

	if err = photoMgr.Update(photo); err != nil {
		panic(err)
	}

	user.RegisterVote(photo.ID)

	if err = userMgr.Update(user); err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
}
