package models

import (
    "testing"
    "strconv"
)

func TestGetPhotoIfNotNone(t *testing.T) {

	tdb := MakeTestDB()
	defer tdb.Clean()

	user := &User{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := user.Insert(); err != nil {
		panic(err)
	}
	photo := &Photo{Title: "test", OwnerID: user.ID, Photo: "test.jpg"}
	if err := photo.Insert(); err != nil {
		panic(err)
	}

	photo, err := GetPhoto(strconv.FormatInt(photo.ID, 10))
	if err != nil {
		panic(err)
	}
	if photo == nil {
		t.Error("Photo should not be nil")
	}

}

func TestGetPhotoIfNone(t *testing.T) {

	tdb := MakeTestDB()
	defer tdb.Clean()

	photo, err := GetPhoto("1")
	if err != nil {
		panic(err)
	}
	if photo != nil {
		t.Error("Photo should be nil")
	}

}

func TestGetAllPhotos(t *testing.T) {
	tdb := MakeTestDB()
	defer tdb.Clean()

	user := &User{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := user.Insert(); err != nil {
		panic(err)
	}
	photo := &Photo{Title: "test", OwnerID: user.ID, Photo: "test.jpg"}
	if err := photo.Insert(); err != nil {
		panic(err)
	}
	photos, err := GetPhotos(1)
	if err != nil {
		panic(err)
	}

	if len(photos) != 1 {
		t.Error("There should be 1 photo")
	}
}
