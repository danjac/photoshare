package models

import (
	"github.com/coopernurse/gorp"
	"image"
	"image/jpeg"
	"image/png"
    "mime/multipart"
	"github.com/nfnt/resize"
	"io"
    "os"
    "time"
    "strings"
)

const (
	UploadsDir = "app/uploads"
)


type Photo struct {
	ID        int       `db:"id" json:"id"`
	OwnerID   int       `db:"owner_id" json:"ownerId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Title     string    `db:"title" json:"title"`
	Photo     string    `db:"photo" json:"photo"`
	Thumbnail string    `db:"thumbnail" json:"thumbnail"`
}

func (photo *Photo) PreInsert(s gorp.SqlExecutor) error {
	photo.CreatedAt = time.Now()
	return nil
}

func (photo *Photo) ProcessImage(src multipart.File, filename, contentType string) error {
	if err := os.MkdirAll(UploadsDir+"/thumbnails", 0777); err != nil && !os.IsExist(err) {
		return err
	}

	// make thumbnail
	var (
		img image.Image
		err error
	)

	if contentType == "image/png" {
		img, err = png.Decode(src)
	} else {
		img, err = jpeg.Decode(src)
	}

	if err != nil {
		return err
	}

	thumb := resize.Thumbnail(300, 300, img, resize.Lanczos3)
	dst, err := os.Create(strings.Join([]string{UploadsDir, "thumbnails", filename}, "/"))

	if err != nil {
		return err
	}

	defer dst.Close()

	if contentType == "image/png" {
		png.Encode(dst, thumb)
	} else if contentType == "image/jpeg" {
		jpeg.Encode(dst, thumb, nil)
	}

	src.Seek(0, 0)

	dst, err = os.Create(strings.Join([]string{UploadsDir, filename}, "/"))

	if err != nil {
		return err
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	photo.Photo = filename
	if _, err := dbMap.Update(photo); err != nil {
		return err
	}

	return nil
}

func (photo *Photo) Save() error {
    return dbMap.Insert(photo)
}

func GetPhotos() ([]Photo, error) {
	var photos []Photo
	if _, err := dbMap.Select(&photos, "SELECT * FROM photos WHERE photo != '' AND photo IS NOT NULL  ORDER BY created_at DESC"); err != nil {
		return photos, err
	}
    return photos, nil
}
