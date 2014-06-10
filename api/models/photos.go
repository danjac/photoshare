package models

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	"github.com/danjac/photoshare/api/utils"
	"os"
	"strings"
	"time"
)

const (
	PageSize   = 32
	UploadsDir = "./public/uploads"
)

func GetImageProcessor() utils.ImageProcessor {
	return utils.LocalImageProcessor{UploadsDir, "/thumbnails"}
}

type Photo struct {
	ID        int64     `db:"id" json:"id"`
	OwnerID   int64     `db:"owner_id" json:"ownerId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Title     string    `db:"title" json:"title"`
	Photo     string    `db:"photo" json:"photo"`
}

func (photo *Photo) PreInsert(s gorp.SqlExecutor) error {
	photo.CreatedAt = time.Now()
	return nil
}

func (photo *Photo) PreDelete(s gorp.SqlExecutor) error {
	filename := strings.Join([]string{UploadsDir, photo.Photo}, "/")
	if err := os.Remove(filename); err != nil {
		return err
	}
	thumbnail := strings.Join([]string{UploadsDir, "thumbnails", photo.Photo}, "/")
	if err := os.Remove(thumbnail); err != nil {
		return err
	}
	return nil
}

func (photo *Photo) CanDelete(user *User) bool {
	return user.ID == photo.OwnerID || user.IsAdmin
}

func (photo *Photo) CanEdit(user *User) bool {
	return user.ID == photo.OwnerID
}

func (photo *Photo) Delete() error {
	_, err := dbMap.Delete(photo)
	return err
}

func (photo *Photo) Update() error {
	_, err := dbMap.Update(photo)
	return err
}

func (photo *Photo) Insert() error {
	return dbMap.Insert(photo)
}

func (photo *Photo) Validate() *ValidationResult {
	result := NewValidationResult()
	if photo.OwnerID == 0 {
		result.Error("owner_id", "Owner ID is missing")
	}
	if photo.Title == "" {
		result.Error("title", "Title is missing")
	}
	if len(photo.Title) > 200 {
		result.Error("title", "Title is too long")
	}
	if photo.Photo == "" {
		result.Error("photo", "Photo filename not set")
	}
	return result
}

type PhotoDetail struct {
	ID        int64     `db:"id" json:"id"`
	OwnerID   int64     `db:"owner_id" json:"ownerId"`
	OwnerName string    `db:"owner_name" json:"ownerName"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Title     string    `db:"title" json:"title"`
	Photo     string    `db:"photo" json:"photo"`
}

func GetPhoto(photoID string) (*Photo, error) {

	photo := &Photo{}
	obj, err := dbMap.Get(photo, photoID)
	if err != nil {
		return photo, err
	}
	if obj == nil {
		return nil, nil
	}
	return obj.(*Photo), nil
}

func GetPhotoDetail(photoID string) (*PhotoDetail, error) {

	photo := &PhotoDetail{}

	q := "SELECT p.*, u.name AS owner_name " +
		"FROM photos p JOIN users u ON u.id = p.owner_id " +
		"WHERE p.id=$1"

	if err := dbMap.SelectOne(photo, q, photoID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return photo, nil

}

func getOffset(pageNum int64) int64 {
	return (pageNum - 1) * PageSize
}

func SearchPhotos(pageNum int64, q string) ([]Photo, error) {

	var photos []Photo
	offset := getOffset(pageNum)

	q = "%" + q + "%"
	if _, err := dbMap.Select(&photos,
		"SELECT * FROM photos WHERE title ILIKE $1 "+
			"ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		q, PageSize, offset); err != nil {
		return photos, err
	}
	return photos, nil

}

func GetPhotos(pageNum int64) ([]Photo, error) {

	var photos []Photo

	offset := getOffset(pageNum)

	if _, err := dbMap.Select(&photos, "SELECT * FROM photos ORDER BY created_at DESC LIMIT $1 OFFSET $2", PageSize, offset); err != nil {
		return photos, err
	}
	return photos, nil
}
