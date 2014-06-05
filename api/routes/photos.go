package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/render"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/utils"
	"net/http"
    "strconv"
)

func upload(w http.ResponseWriter, r *http.Request) {

	user, err := session.GetCurrentUser(r)
	if err != nil {
		render.Error(w, err)
		return
	}
	if user == nil {
		render.Status(w, http.StatusUnauthorized, "Not logged in")
		return
	}

	title := r.FormValue("title")
	src, hdr, err := r.FormFile("photo")
	if err != nil {
		render.Error(w, err)
		return
	}
	contentType := hdr.Header["Content-Type"][0]
	if contentType != "image/png" && contentType != "image/jpeg" {
		render.Status(w, http.StatusBadRequest, "Not a valid image")
		return
	}

	defer src.Close()
	filename, err := utils.ProcessImage(src, contentType, fileUploadDir)
	if err != nil {
		render.Error(w, err)
		return
	}

	photo := &models.Photo{Title: title,
		OwnerID: user.ID, Photo: filename}

	if result := photo.Validate(); !result.OK {
		render.JSON(w, http.StatusBadRequest, result)
		return
	}

	if err := photo.Save(); err != nil {
		render.Error(w, err)
		return
	}

	render.JSON(w, http.StatusOK, photo)
}

func getPhotos(w http.ResponseWriter, r *http.Request) {

    pageNum, err := strconv.ParseInt(r.FormValue("page"), 10, 0)
    if err != nil {
        pageNum = 1
    }

	photos, err := models.GetPhotos(pageNum)
	if err != nil {
		render.Error(w, err)
		return
	}
	render.JSON(w, http.StatusOK, photos)
}
