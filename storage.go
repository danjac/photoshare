package photoshare

import (
	"code.google.com/p/graphics-go/graphics"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/disintegration/gift"
	"github.com/juju/errgo"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
)

const (
	thumbnailHeight = 300
	thumbnailWidth  = 300
)

type readable interface {
	io.Reader
	io.Seeker
}

var allowedContentTypes = []string{
	"image/png",
	"image/jpeg",
	"image/gif"}

func isAllowedContentType(contentType string) bool {
	fmt.Println("CHECKING", contentType)
	for _, value := range allowedContentTypes {
		if contentType == value {
			return true
		}
	}

	return false
}

func generateRandomFilename(contentType string) string {

	var ext string

	switch contentType {
	case "image/jpeg":

		ext = ".jpg"
	case "image/png":
		ext = ".png"

	case "image/gif":
		ext = ".gif"
	}

	return uniuri.New() + ext
}

type fileStorage interface {
	clean(string) error
	store(readable, string, string) error
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

func (f *defaultFileStorage) store(src readable, filename, contentType string) error {
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

	switch contentType {
	case "image/png":
		img, err = png.Decode(src)
	case "image/jpeg":
		img, err = jpeg.Decode(src)
	case "image/jpg":
		img, err = jpeg.Decode(src)
	case "image/gif":
		img, err = gif.Decode(src)
	default:
		return errors.New("invalid content type:" + contentType)
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

	switch contentType {
	case "image/png":
		png.Encode(dst, thumb)
	case "image/jpeg":
		jpeg.Encode(dst, thumb, nil)
	case "image/jpg":
		jpeg.Encode(dst, thumb, nil)
	case "image/gif":
		gif.Encode(dst, thumb, nil)
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
