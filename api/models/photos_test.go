package models

import (
	"strconv"
	"testing"
)

func TestGetIfNotNone(t *testing.T) {

	tdb := MakeTestDB()
	defer tdb.Clean()

	user := &User{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := userMgr.Insert(user); err != nil {
		t.Error(err)
		return
	}
	photo := &Photo{Title: "test", OwnerID: user.ID, Photo: "test.jpg"}
	if err := photoMgr.Insert(photo); err != nil {
		t.Error(err)
		return
	}

	photo, err := NewPhotoManager().Get(strconv.FormatInt(photo.ID, 10))
	if err != nil {
		t.Error(err)
		return
	}
	if photo == nil {
		t.Error("Photo should not be nil")
	}

}

func TestGetIfNone(t *testing.T) {

	tdb := MakeTestDB()
	defer tdb.Clean()

	photo, err := NewPhotoManager().Get("1")
	if err != nil {
		t.Error(err)
		return
	}
	if photo != nil {
		t.Error("Photo should be nil")
	}

}

func TestSearchPhotos(t *testing.T) {
	tdb := MakeTestDB()
	defer tdb.Clean()

	photoMgr := NewPhotoManager()
	userMgr := NewUserManager()

	user := &User{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := userMgr.Insert(user); err != nil {
		t.Error(err)
		return
	}
	photo := &Photo{Title: "test", OwnerID: user.ID, Photo: "test.jpg"}
	if err := photoMgr.Insert(photo); err != nil {
		t.Error(err)
		return
	}
	photos, err := photoMgr.Search(1, "test")
	if err != nil {
		t.Error(err)
		return
	}

	if len(photos) != 1 {
		t.Error("There should be 1 photo")
	}
}
func TestGetPhotos(t *testing.T) {
	tdb := MakeTestDB()
	defer tdb.Clean()

	photoMgr := NewPhotoManager()
	userMgr := NewUserManager()

	user := &User{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := userMgr.Insert(user); err != nil {
		t.Error(err)
		return
	}
	photo := &Photo{Title: "test", OwnerID: user.ID, Photo: "test.jpg"}
	if err := photoMgr.Insert(photo); err != nil {
		t.Error(err)
		return
	}
	photos, err := photoMgr.All(1)
	if err != nil {
		t.Error(err)
		return
	}

	if len(photos) != 1 {
		t.Error("There should be 1 photo")
	}
}
