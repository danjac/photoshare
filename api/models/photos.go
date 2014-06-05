package models

import (
	"github.com/coopernurse/gorp"
	"time"
)

type Photo struct {
	ID        int       `db:"id" json:"id"`
	OwnerID   int       `db:"owner_id" json:"ownerId"`
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

func GetPhotos() ([]Photo, error) {
	var photos []Photo
	if _, err := dbMap.Select(&photos, "SELECT * FROM photos WHERE photo != '' AND photo IS NOT NULL  ORDER BY created_at DESC"); err != nil {
		return photos, err
	}
	return photos, nil
}
