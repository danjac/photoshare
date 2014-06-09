package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func deletePhoto(w http.ResponseWriter, r *http.Request) {

	user, err := session.GetCurrentUser(r)
	if err != nil {
		panic(err)
	}

	if user == nil {
		render(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	photo, err := models.GetPhoto(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}
	if photo == nil {
		render(w, http.StatusNotFound, "Photo not found")
		return
	}

	if !photo.CanDelete(user) {
		render(w, http.StatusForbidden, "You can't delete this photo")
		return
	}
	if err := photo.Delete(); err != nil {
		panic(err)
	}

	render(w, http.StatusOK, "Photo deleted")
}

func photoDetail(w http.ResponseWriter, r *http.Request) {

	photo, err := models.GetPhotoDetail(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}
	if photo == nil {
		render(w, http.StatusNotFound, "Photo not found")
		return
	}

	render(w, http.StatusOK, photo)
}

func editPhoto(w http.ResponseWriter, r *http.Request) {

	user, err := session.GetCurrentUser(r)
	if err != nil {
		panic(err)
	}

	if user == nil {
		render(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	photo, err := models.GetPhoto(mux.Vars(r)["id"])
	if err != nil {
		panic(err)
	}

	if photo == nil {
		render(w, http.StatusNotFound, "No photo found")
		return
	}

	if !photo.CanEdit(user) {
		render(w, http.StatusForbidden, "You can't edit this photo")
		return
	}

	newPhoto := &models.Photo{}

	if err := parseJSON(r, newPhoto); err != nil {
		panic(err)
	}

	photo.Title = newPhoto.Title

	if result := photo.Validate(); !result.OK {
		render(w, http.StatusBadRequest, result)
		return
	}

	if err := photo.Update(); err != nil {
		panic(err)
	}

	render(w, http.StatusOK, photo)
}

func upload(w http.ResponseWriter, r *http.Request) {

	user, err := session.GetCurrentUser(r)
	if err != nil {
		panic(err)
	}

	if user == nil {
		render(w, http.StatusUnauthorized, "You must be logged in")
		return
	}

	title := r.FormValue("title")
	src, hdr, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile {
			render(w, http.StatusBadRequest, "No image was posted")
			return
		}
		panic(err)
	}
	contentType := hdr.Header["Content-Type"][0]
	if contentType != "image/png" && contentType != "image/jpeg" {
		render(w, http.StatusBadRequest, "Not a valid image")
		return
	}

	defer src.Close()
	filename, err := utils.ProcessImage(src, contentType)
	if err != nil {
		panic(err)
	}

	photo := &models.Photo{Title: title,
		OwnerID: user.ID, Photo: filename}

	if result := photo.Validate(); !result.OK {
		render(w, http.StatusBadRequest, result)
		return
	}

	if err := photo.Insert(); err != nil {
		panic(err)
	}

	render(w, http.StatusOK, photo)
}

func getPhotos(w http.ResponseWriter, r *http.Request) {

	pageNum, err := strconv.ParseInt(r.FormValue("page"), 10, 0)
	if err != nil {
		pageNum = 1
	}

	photos, err := models.GetPhotos(pageNum)
	if err != nil {
		panic(err)
	}
	render(w, http.StatusOK, photos)
}
