package api

import (
	"code.google.com/p/graphics-go/graphics"
	"errors"
	"github.com/dchest/uniuri"
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

type FileManager interface {
	Clean(string) error
	Store(src multipart.File, contentType string) (string, error)
}

func NewFileManager(config *AppConfig) FileManager {
	return &defaultFileManager{config}
}

type defaultFileManager struct {
	config *AppConfig
}

func (f *defaultFileManager) Clean(name string) error {

	imagePath := path.Join(f.config.UploadsDir, name)
	thumbnailPath := path.Join(f.config.ThumbnailsDir, name)

	if err := os.Remove(imagePath); err != nil {
		return err
	}
	if err := os.Remove(thumbnailPath); err != nil {
		return err
	}
	return nil
}

func (f *defaultFileManager) Store(src multipart.File, contentType string) (string, error) {

	if !isAllowedContentType(contentType) {
		return "", InvalidContentType
	}
	filename := generateRandomFilename(contentType)

	if err := os.MkdirAll(f.config.UploadsDir, 0777); err != nil && !os.IsExist(err) {
		return filename, err
	}

	if err := os.MkdirAll(f.config.ThumbnailsDir, 0777); err != nil && !os.IsExist(err) {
		return filename, err
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
		return filename, err
	}

	thumb := image.NewRGBA(image.Rect(0, 0, ThumbnailWidth, ThumbnailHeight))
	graphics.Thumbnail(thumb, img)

	dst, err := os.Create(path.Join(f.config.ThumbnailsDir, filename))

	if err != nil {
		return filename, err
	}

	defer dst.Close()

	if contentType == "image/png" {
		png.Encode(dst, thumb)
	} else if contentType == "image/jpeg" {
		jpeg.Encode(dst, thumb, nil)
	}

	src.Seek(0, 0)

	dst, err = os.Create(path.Join(f.config.UploadsDir, filename))

	if err != nil {
		return filename, err
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return filename, err
	}

	return filename, nil

}
