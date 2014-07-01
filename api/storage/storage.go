package storage

import (
	"github.com/danjac/photoshare/api/config"
	"os"
	"path"
)

type PhotoCleaner interface {
	Clean(string) error
}

type defaultPhotoCleaner struct {
}

func (c *defaultPhotoCleaner) Clean(name string) error {

	imagePath := path.Join(config.Dirs.Uploads, name)
	thumbnailPath := path.Join(config.Dirs.Thumbnails, name)

	if err := os.Remove(imagePath); err != nil {
		return err
	}
	if err := os.Remove(thumbnailPath); err != nil {
		return err
	}
	return nil
}

var photoCleaner = &defaultPhotoCleaner{}

func NewPhotoCleaner() PhotoCleaner {
	return photoCleaner
}
