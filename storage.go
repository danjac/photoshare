package photoshare

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
	thumbnailHeight = 300
	thumbnailWidth  = 300
)

var (
	allowedContentTypes   = []string{"image/png", "image/jpeg"}
	errInvalidContentType = errors.New("must be PNG or JPG")
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

type fileStorage interface {
	clean(string) error
	store(src multipart.File, contentType string) (string, error)
}

func newFileStorage(cfg *config) fileStorage {
	return &defaultFileStorage{
		cfg.UploadsDir,
		cfg.ThumbnailsDir,
	}
}

type defaultFileStorage struct {
	uploadsDir, thumbnailsDir string
}

func (f *defaultFileStorage) clean(name string) error {

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

func (f *defaultFileStorage) saveFiles(filename string, contentType string, src multipart.File) error {
	if err := os.MkdirAll(f.uploadsDir, 0777); err != nil && !os.IsExist(err) {
		return errgo.Mask(err)
	}

	if err := os.MkdirAll(f.thumbnailsDir, 0777); err != nil && !os.IsExist(err) {
		return errgo.Mask(err)
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
		return errgo.Mask(err)
	}

	thumb := image.NewRGBA(image.Rect(0, 0, thumbnailWidth, thumbnailHeight))
	graphics.Thumbnail(thumb, img)

	dst, err := os.Create(path.Join(f.thumbnailsDir, filename))
	if err != nil {
		return errgo.Mask(err)
	}

	g := gift.New(gift.Contrast(-30))
	g.Draw(thumb, thumb)

	if err != nil {
		return errgo.Mask(err)
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
		return errgo.Mask(err)
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return errgo.Mask(err)
	}

	return nil

}

func (f *defaultFileStorage) store(src multipart.File, contentType string) (string, error) {

	if !isAllowedContentType(contentType) {
		return "", errInvalidContentType
	}

	filename := generateRandomFilename(contentType)

	go func() {
		err := f.saveFiles(filename, contentType, src)
		if err != nil {
			logError(err)
		}
	}()

	return filename, nil

}
