package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/storage"
	"github.com/danjac/photoshare/api/validation"
	"net/http"
	"strconv"
	"strings"
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

func deletePhoto(c *Context) *Result {

	photo, err := photoMgr.Get(c.Param("id"))
	if err != nil {
		return c.Error(err)
	}
	if photo == nil {
		return c.NotFound("Photo not found")
	}

	if !photo.CanDelete(c.User) {
		return c.Forbidden("You can't delete this photo")
	}
	if err := photoMgr.Delete(photo); err != nil {
		return c.Error(err)
	}

	return c.OK("Photo deleted")
}

func photoDetail(c *Context) *Result {

	user, err := c.GetCurrentUser()
	if err != nil {
		return c.Error(err)
	}
	photo, err := photoMgr.GetDetail(c.Param("id"), user)
	if err != nil {
		return c.Error(err)
	}
	if photo == nil {
		return c.NotFound("Photo not found")
	}

	return c.OK(photo)
}

func getPhotoToEdit(c *Context) (*models.Photo, *Result) {
	photo, err := photoMgr.Get(c.Param("id"))
	if err != nil {
		return photo, c.Error(err)
	}

	if photo == nil {
		return photo, c.NotFound("No photo found")
	}

	if !photo.CanEdit(c.User) {
		return photo, c.Forbidden("You can't edit this photo")
	}
	return photo, nil
}

func editPhotoTitle(c *Context) *Result {

	photo, result := getPhotoToEdit(c)

	if result != nil {
		return result
	}

	s := &struct {
		Title string `json:"title"`
	}{}

	if err := c.ParseJSON(s); err != nil {
		return c.Error(err)
	}

	photo.Title = s.Title

	validator := validation.NewPhotoValidator(photo)

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			return c.Error(err)
		}
		return c.BadRequest(result)
	}

	if err := photoMgr.Update(photo); err != nil {
		return c.Error(err)
	}

	return c.OK("Photo Updated")
}

func editPhotoTags(c *Context) *Result {

	photo, result := getPhotoToEdit(c)

	if result != nil {
		return result
	}

	s := &struct {
		Tags []string `json:"tags"`
	}{}

	if err := c.ParseJSON(s); err != nil {
		return c.Error(err)
	}

	photo.Tags = s.Tags

	if err := photoMgr.UpdateTags(photo); err != nil {
		return c.Error(err)
	}

	return c.OK("Photo updated")
}

func upload(c *Context) *Result {

	title := c.FormValue("title")
	taglist := c.FormValue("taglist")
	tags := strings.Split(taglist, " ")

	src, hdr, err := c.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			return c.BadRequest("No image was posted")
		}
		return c.Error(err)
	}
	contentType := hdr.Header["Content-Type"][0]

	if !isAllowedContentType(contentType) {
		return c.BadRequest("Not a valid image")
	}

	defer src.Close()

	processor := storage.NewImageProcessor()
	filename, err := processor.Process(src, contentType)

	if err != nil {
		return c.Error(err)
	}

	photo := &models.Photo{Title: title,
		OwnerID: c.User.ID, Filename: filename, Tags: tags}

	validator := validation.NewPhotoValidator(photo)

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			return c.Error(err)
		}
		return c.BadRequest(result)
	}

	if err := photoMgr.Insert(photo); err != nil {
		return c.Error(err)
	}

	return c.OK(photo)
}

func getPhotos(c *Context) *Result {
	var (
		err    error
		photos []models.Photo
	)

	pageNum, err := strconv.ParseInt(c.FormValue("page"), 10, 64)
	if err != nil {
		pageNum = 1
	}

	q := c.FormValue("q")
	ownerID := c.FormValue("ownerID")
	orderBy := c.FormValue("orderBy")

	if q != "" {
		photos, err = photoMgr.Search(pageNum, q)
	} else if ownerID != "" {
		photos, err = photoMgr.ByOwnerID(pageNum, ownerID)
	} else {
		photos, err = photoMgr.All(pageNum, orderBy)
	}

	if err != nil {
		return c.Error(err)
	}
	return c.OK(photos)
}

func getTags(c *Context) *Result {
	tags, err := photoMgr.GetTagCounts()
	if err != nil {
		return c.Error(err)
	}
	return c.OK(tags)
}

func voteDown(c *Context) *Result {
	return vote(c, func(photo *models.Photo) { photo.DownVotes += 1 })
}

func voteUp(c *Context) *Result {
	return vote(c, func(photo *models.Photo) { photo.UpVotes += 1 })
}

func vote(c *Context, fn func(photo *models.Photo)) *Result {
	var (
		photo *models.Photo
		err   error
	)

	photo, err = photoMgr.Get(c.Param("id"))
	if err != nil {
		return c.Error(err)
	}
	if photo == nil {
		return c.NotFound("Photo not found")
	}

	if !photo.CanVote(c.User) {
		return c.Forbidden("You can't vote on this photo")
	}

	fn(photo)

	if err = photoMgr.Update(photo); err != nil {
		return c.Error(err)
	}

	c.User.AddVote(photo.ID)

	if err = userMgr.Update(c.User); err != nil {
		return c.Error(err)
	}
	return c.OK("OK")
}
