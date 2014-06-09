package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/utils"
	"net/http"
	"strconv"
)

func deletePhoto(c *AppContext) {

	photo, err := models.GetPhoto(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	if photo == nil {
		c.NotFound("Photo not found")
		return
	}

	if !photo.CanDelete(c.User) {
		c.Forbidden("You can't delete this photo")
		return
	}
	if err := photo.Delete(); err != nil {
		c.Error(err)
		return
	}

	c.OK("Photo deleted")
}

func photoDetail(c *AppContext) {

	photo, err := models.GetPhotoDetail(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}
	if photo == nil {
		c.NotFound("Photo not found")
		return
	}

	c.OK(photo)
}

func editPhoto(c *AppContext) {

	photo, err := models.GetPhoto(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	if photo == nil {
		c.NotFound("No photo found")
		return
	}

	if !photo.CanEdit(c.User) {
		c.Forbidden("You can't edit this photo")
		return
	}

	newPhoto := &models.Photo{}

	if err := c.ParseJSON(newPhoto); err != nil {
		c.Error(err)
		return
	}

	photo.Title = newPhoto.Title

	if result := photo.Validate(); !result.OK {
		c.BadRequest(result)
		return
	}

	if err := photo.Update(); err != nil {
		c.Error(err)
		return
	}

	c.OK(photo)
}

func upload(c *AppContext) {

	title := c.FormValue("title")
	src, hdr, err := c.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile {
			c.BadRequest("No image was posted")
			return
		}
		c.Error(err)
		return
	}
	contentType := hdr.Header["Content-Type"][0]
	if contentType != "image/png" && contentType != "image/jpeg" {
		c.BadRequest("Not a valid image")
		return
	}

	defer src.Close()
	filename, err := utils.ProcessImage(src, contentType)
	if err != nil {
		c.Error(err)
		return
	}

	photo := &models.Photo{Title: title,
		OwnerID: c.User.ID, Photo: filename}

	if result := photo.Validate(); !result.OK {
		c.BadRequest(result)
		return
	}

	if err := photo.Insert(); err != nil {
		c.Error(err)
		return
	}

	c.OK(photo)
}

func getPhotos(c *AppContext) {

	pageNum, err := strconv.ParseInt(c.FormValue("page"), 10, 0)
	if err != nil {
		pageNum = 1
	}

	photos, err := models.GetPhotos(pageNum)
	if err != nil {
		c.Error(err)
		return
	}
	c.OK(photos)
}
