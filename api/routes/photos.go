package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/utils"
	"net/http"
	"strconv"
)

var allowedContentTypes = []string{"image/png", "image/jpeg"}

func isAllowedContentType(contentType string) bool {
	for _, value := range allowedContentTypes {
		if contentType == value {
			return true
		}
	}

	return false
}

func deletePhoto(c *AppContext) error {

	photo, err := models.GetPhoto(c.Param("id"))
	if err != nil {
		return err
	}
	if photo == nil {
		return c.NotFound("Photo not found")
	}

	if !photo.CanDelete(c.User) {
		return c.Forbidden("You can't delete this photo")
	}
	if err := photo.Delete(); err != nil {
		return err
	}

	return c.OK("Photo deleted")
}

func photoDetail(c *AppContext) error {

	photo, err := models.GetPhotoDetail(c.Param("id"))
	if err != nil {
		return err
	}
	if photo == nil {
		return c.NotFound("Photo not found")
	}

	return c.OK(photo)
}

func editPhoto(c *AppContext) error {

	photo, err := models.GetPhoto(c.Param("id"))
	if err != nil {
		return err
	}

	if photo == nil {
		return c.NotFound("No photo found")
	}

	if !photo.CanEdit(c.User) {
		return c.Forbidden("You can't edit this photo")
	}

	newPhoto := &models.Photo{}

	if err := c.ParseJSON(newPhoto); err != nil {
		return err
	}

	photo.Title = newPhoto.Title

	if result := photo.Validate(); !result.OK {
		return c.BadRequest(result)
	}

	if err := photo.Update(); err != nil {
		return err
	}

	return c.OK(photo)
}

func upload(c *AppContext) error {

	title := c.FormValue("title")
	src, hdr, err := c.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile {
			return c.BadRequest("No image was posted")
		}
		return err
	}
	contentType := hdr.Header["Content-Type"][0]

	if !isAllowedContentType(contentType) {
		return c.BadRequest("Not a valid image")
	}

	defer src.Close()
	filename, err := utils.ProcessImage(src, contentType)
	if err != nil {
		return err
	}

	photo := &models.Photo{Title: title,
		OwnerID: c.User.ID, Photo: filename}

	if result := photo.Validate(); !result.OK {
		return c.BadRequest(result)
	}

	if err := photo.Insert(); err != nil {
		return err
	}

	return c.OK(photo)
}

func getPhotos(c *AppContext) error {
	var (
		err    error
		photos []models.Photo
	)

	pageNum, err := strconv.ParseInt(c.FormValue("page"), 10, 64)
	if err != nil {
		pageNum = 1
	}

	q := c.FormValue("q")

	if q == "" {
		photos, err = models.GetPhotos(pageNum)
	} else {
		photos, err = models.SearchPhotos(pageNum, q)
	}
	if err != nil {
		return err
	}
	return c.OK(photos)
}
