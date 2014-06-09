package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/utils"
	"net/http"
	"strconv"
)

func deletePhoto(c *AppContext) {

	user, err := c.GetCurrentUser()
	if err != nil {
		c.Error(err)
		return
	}

	if user == nil {
		c.Render(http.StatusUnauthorized, "You must be logged in")
		return
	}

	photo, err := models.GetPhoto(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	if photo == nil {
		c.Render(http.StatusNotFound, "Photo not found")
		return
	}

	if !photo.CanDelete(user) {
		c.Render(http.StatusForbidden, "You can't delete this photo")
		return
	}
	if err := photo.Delete(); err != nil {
		c.Error(err)
		return
	}

	c.Render(http.StatusOK, "Photo deleted")
}

func photoDetail(c *AppContext) {

	photo, err := models.GetPhotoDetail(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	if photo == nil {
		c.Render(http.StatusNotFound, "Photo not found")
		return
	}

	c.Render(http.StatusOK, photo)
}

func editPhoto(c *AppContext) {

	user, err := c.GetCurrentUser()
	if err != nil {
		c.Error(err)
		return
	}

	if user == nil {
		c.Render(http.StatusUnauthorized, "You must be logged in")
		return
	}

	photo, err := models.GetPhoto(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	if photo == nil {
		c.Render(http.StatusNotFound, "No photo found")
		return
	}

	if !photo.CanEdit(user) {
		c.Render(http.StatusForbidden, "You can't edit this photo")
		return
	}

	newPhoto := &models.Photo{}

	if err := c.ParseJSON(newPhoto); err != nil {
		c.Error(err)
		return
	}

	photo.Title = newPhoto.Title

	if result := photo.Validate(); !result.OK {
		c.Render(http.StatusBadRequest, result)
		return
	}

	if err := photo.Update(); err != nil {
		c.Error(err)
		return
	}

	c.Render(http.StatusOK, photo)
}

func upload(c *AppContext) {

	user, err := c.GetCurrentUser()
	if err != nil {
		c.Error(err)
		return
	}

	if user == nil {
		c.Render(http.StatusUnauthorized, "You must be logged in")
		return
	}

	title := c.Request.FormValue("title")
	src, hdr, err := c.Request.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile {
			c.Render(http.StatusBadRequest, "No image was posted")
			return
		}
		c.Error(err)
		return
	}
	contentType := hdr.Header["Content-Type"][0]
	if contentType != "image/png" && contentType != "image/jpeg" {
		c.Render(http.StatusBadRequest, "Not a valid image")
		return
	}

	defer src.Close()
	filename, err := utils.ProcessImage(src, contentType)
	if err != nil {
		c.Error(err)
		return
	}

	photo := &models.Photo{Title: title,
		OwnerID: user.ID, Photo: filename}

	if result := photo.Validate(); !result.OK {
		c.Render(http.StatusBadRequest, result)
		return
	}

	if err := photo.Insert(); err != nil {
		c.Error(err)
		return
	}

	c.Render(http.StatusOK, photo)
}

func getPhotos(c *AppContext) {

	pageNum, err := strconv.ParseInt(c.Request.FormValue("page"), 10, 0)
	if err != nil {
		pageNum = 1
	}

	photos, err := models.GetPhotos(pageNum)
	if err != nil {
		c.Error(err)
		return
	}
	c.Render(http.StatusOK, photos)
}
