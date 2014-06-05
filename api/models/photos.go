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

func GetPhotos() ([]Photo, error) {
	var photos []Photo
	if _, err := dbMap.Select(&photos, "SELECT * FROM photos WHERE photo != '' AND photo IS NOT NULL  ORDER BY created_at DESC"); err != nil {
		return photos, err
	}
	return photos, nil
}
