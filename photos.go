package photoshare

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func deletePhoto(ctx *context, w http.ResponseWriter, r *http.Request) error {

	photo, err := ctx.datamapper.getPhoto(ctx.params.getInt("id"))
	if err != nil {
		return err
	}

	if !photo.canDelete(ctx.user) {
		return httpError{http.StatusForbidden, "You're not allowed to delete this photo"}
	}
	if err := ctx.datamapper.removePhoto(photo); err != nil {
		return err
	}

	go func() {
		if err := ctx.filestore.clean(photo.Filename); err != nil {
			log.Println(err)
		}
	}()

	if err := ctx.cache.clear(); err != nil {
		return err
	}

	sendMessage(&socketMessage{ctx.user.Name, "", photo.ID, "photo_deleted"})
	return renderString(w, http.StatusOK, "Photo deleted")
}

func getPhotoDetail(ctx *context, w http.ResponseWriter, r *http.Request) error {

	photo, err := ctx.datamapper.getPhotoDetail(ctx.params.getInt("id"), ctx.user)
	if err != nil {
		return err
	}
	return renderJSON(w, photo, http.StatusOK)

}

func getPhotoToEdit(ctx *context, w http.ResponseWriter, r *http.Request) (*photo, error) {

	photo, err := ctx.datamapper.getPhoto(ctx.params.getInt("id"))
	if err != nil {
		return photo, err
	}

	if !photo.canEdit(ctx.user) {
		return photo, httpError{http.StatusForbidden, "You're not allowed to edit this photo"}
	}
	return photo, nil
}

func editPhotoTitle(ctx *context, w http.ResponseWriter, r *http.Request) error {

	photo, err := getPhotoToEdit(ctx, w, r)

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

	if err := ctx.validate(photo, r); err != nil {
		return err

	}

	if err := ctx.datamapper.updatePhoto(photo); err != nil {
		return err
	}

	sendMessage(&socketMessage{ctx.user.Name, "", photo.ID, "photo_updated"})
	return renderString(w, http.StatusOK, "Photo updated")
}

func editPhotoTags(ctx *context, w http.ResponseWriter, r *http.Request) error {

	photo, err := getPhotoToEdit(ctx, w, r)
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
	if err := ctx.datamapper.updateTags(photo); err != nil {
		return err
	}

	sendMessage(&socketMessage{ctx.user.Name, "", photo.ID, "photo_updated"})
	return renderString(w, http.StatusOK, "Photo updated")

}

func upload(ctx *context, w http.ResponseWriter, r *http.Request) error {

	title := r.FormValue("title")
	taglist := r.FormValue("taglist")
	tags := strings.Split(taglist, " ")

	src, hdr, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			return httpError{http.StatusBadRequest, "Invalid photo"}
		}
		return err
	}
	defer src.Close()

	contentType := hdr.Header["Content-Type"][0]

	if !isAllowedContentType(contentType) {
		return httpError{http.StatusBadRequest, "Only JPEG or PNG files allowed"}
	}

	filename := generateRandomFilename(contentType)

	photo := &photo{Title: title,
		OwnerID:  ctx.user.ID,
		Filename: filename,
		Tags:     tags,
	}

	go func() {
		err := ctx.filestore.store(src, photo.Filename, contentType)
		if err != nil {
			logError(err)
		}
	}()

	if err := ctx.validate(photo, r); err != nil {
		return err
	}
	if err := ctx.datamapper.createPhoto(photo); err != nil {
		return err
	}
	if err := ctx.cache.clear(); err != nil {
		logError(err)
	}

	sendMessage(&socketMessage{ctx.user.Name, "", photo.ID, "photo_uploaded"})
	return renderJSON(w, photo, http.StatusCreated)
}

func searchPhotos(ctx *context, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	q := r.FormValue("q")
	cacheKey := fmt.Sprintf("photos:search:%s:page:%d", q, page.index)

	return ctx.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := ctx.datamapper.searchPhotos(page, q)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})

}

func photosByOwnerID(ctx *context, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	ownerID := ctx.params.getInt("ownerID")
	cacheKey := fmt.Sprintf("photos:ownerID:%d:page:%d", ownerID, page.index)

	return ctx.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := ctx.datamapper.getPhotosByOwnerID(page, ownerID)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func getPhotos(ctx *context, w http.ResponseWriter, r *http.Request) error {

	page := getPage(r)
	orderBy := r.FormValue("orderBy")
	cacheKey := fmt.Sprintf("photos:%s:page:%d", orderBy, page.index)

	return ctx.cache.render(w, http.StatusOK, cacheKey, func() (interface{}, error) {
		photos, err := ctx.datamapper.getPhotos(page, orderBy)
		if err != nil {
			return photos, err
		}
		return photos, nil
	})
}

func getTags(ctx *context, w http.ResponseWriter, r *http.Request) error {
	return ctx.cache.render(w, http.StatusOK, "tags", func() (interface{}, error) {
		tags, err := ctx.datamapper.getTagCounts()
		if err != nil {
			return tags, err
		}
		return tags, nil
	})

}

func voteDown(ctx *context, w http.ResponseWriter, r *http.Request) error {
	return vote(ctx, w, r, func(photo *photo) { photo.DownVotes++ })
}

func voteUp(ctx *context, w http.ResponseWriter, r *http.Request) error {
	return vote(ctx, w, r, func(photo *photo) { photo.UpVotes++ })
}

func vote(ctx *context, w http.ResponseWriter, r *http.Request, fn func(photo *photo)) error {

	photo, err := ctx.datamapper.getPhoto(ctx.params.getInt("id"))
	if err != nil {
		return err
	}

	if !photo.canVote(ctx.user) {
		return httpError{http.StatusForbidden, "You're not allowed to vote on this photo"}
	}

	fn(photo)

	ctx.user.registerVote(photo.ID)

	if err := ctx.datamapper.updateMany(photo, ctx.user); err != nil {
		return err
	}

	return renderString(w, http.StatusOK, "Voting successful")
}
