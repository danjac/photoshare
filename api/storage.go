package api

import (
	"code.google.com/p/graphics-go/graphics"
	"errors"
	"github.com/dchest/uniuri"
	"github.com/disintegration/gift"
	"github.com/juju/errgo"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path"
)

const (
	ThumbnailHeight = 300
	ThumbnailWidth  = 300
)

var (
	allowedContentTypes = []string{"image/png", "image/jpeg"}
	InvalidContentType  = errors.New("Must be PNG or JPG")
)

func isAllowedContentType(contentType string) bool {
	for _, value := range allowedContentTypes {
		if contentType == value {
			return true
		}
	}

	return false
}

func generateRandomFilename(contentType string) string {
	filename := uniuri.New()
	if contentType == "image/png" {
		return filename + ".png"
	}
	return filename + ".jpg"
}

type FileStorage interface {
	Clean(string) error
	Store(src multipart.File, contentType string) (string, error)
}

func NewFileStorage(config *AppConfig) FileStorage {
	return &defaultFileStorage{
		config.UploadsDir,
		config.ThumbnailsDir,
	}
}

type defaultFileStorage struct {
	uploadsDir, thumbnailsDir string
}

func (f *defaultFileStorage) Clean(name string) error {

	imagePath := path.Join(f.uploadsDir, name)
	thumbnailPath := path.Join(f.thumbnailsDir, name)

	if err := os.Remove(imagePath); err != nil {
		return errgo.Mask(err)
	}
	if err := os.Remove(thumbnailPath); err != nil {
		return errgo.Mask(err)
	}
	return nil
}

func (f *defaultFileStorage) Store(src multipart.File, contentType string) (string, error) {

	if !isAllowedContentType(contentType) {
		return "", InvalidContentType
	}
	filename := generateRandomFilename(contentType)

	if err := os.MkdirAll(f.uploadsDir, 0777); err != nil && !os.IsExist(err) {
		return filename, errgo.Mask(err)
	}

	if err := os.MkdirAll(f.thumbnailsDir, 0777); err != nil && !os.IsExist(err) {
		return filename, errgo.Mask(err)
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
		return filename, errgo.Mask(err)
	}


	thumb := image.NewRGBA(image.Rect(0, 0, ThumbnailWidth, ThumbnailHeight))
	graphics.Thumbnail(thumb, img)

	dst, err := os.Create(path.Join(f.thumbnailsDir, filename))

	g := gift.New(gift.Contrast(-30))
	g.Draw(thumb, thumb)

	if err != nil {
		return filename, errgo.Mask(err)
	}

	defer dst.Close()

	if contentType == "image/png" {
		png.Encode(dst, thumb)
	} else if contentType == "image/jpeg" {
		jpeg.Encode(dst, thumb, nil)
	}

	src.Seek(0, 0)

	dst, err = os.Create(path.Join(f.uploadsDir, filename))

	if err != nil {
		return filename, errgo.Mask(err)
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return filename, errgo.Mask(err)
	}

	return filename, nil

}
