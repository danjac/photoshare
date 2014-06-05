package models

import (
	"github.com/coopernurse/gorp"
	"time"
)

const pageSize = 8

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

func (photo *Photo) Save() error {
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

func GetPhotos(pageNum int64) ([]Photo, error) {

    var photos []Photo

    offset := (pageNum - 1) * pageSize

	if _, err := dbMap.Select(&photos, "SELECT * FROM photos ORDER BY created_at DESC LIMIT $1 OFFSET $2", pageSize, offset); err != nil {
		return photos, err
	}
	return photos, nil
}
