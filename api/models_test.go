package api

import (
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
	photo := &Photo{Title: "test", OwnerID: user.ID, Filename: "test.jpg"}
	if err := photoMgr.Insert(photo); err != nil {
		t.Error(err)
		return
	}

	photo, exists, err := NewPhotoManager().Get(photo.ID)
	if err != nil {
		t.Error(err)
		return
	}
	if !exists {
		t.Error("Photo should exist")
	}

}

func TestGetIfNone(t *testing.T) {

	tdb := MakeTestDB()
	defer tdb.Clean()

	_, exists, err := NewPhotoManager().Get(1)
	if err != nil {
		t.Error(err)
		return
	}
	if exists {
		t.Error("Photo should not exist")
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
	photo := &Photo{Title: "test", OwnerID: user.ID, Filename: "test.jpg"}
	if err := photoMgr.Insert(photo); err != nil {
		t.Error(err)
		return
	}
	result, err := photoMgr.Search(1, "test")
	if err != nil {
		t.Error(err)
		return
	}

	if len(result.Items) != 1 {
		t.Error("There should be 1 photo")
	}
}
func TestAllPhotos(t *testing.T) {
	tdb := MakeTestDB()
	defer tdb.Clean()

	photoMgr := NewPhotoManager()
	userMgr := NewUserManager()

	user := &User{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := userMgr.Insert(user); err != nil {
		t.Error(err)
		return
	}
	photo := &Photo{Title: "test", OwnerID: user.ID, Filename: "test.jpg"}
	if err := photoMgr.Insert(photo); err != nil {
		t.Error(err)
		return
	}
	result, err := photoMgr.All(1, "")
	if err != nil {
		t.Error(err)
		return
	}

	if len(result.Items) != 1 {
		t.Error("There should be 1 photo")
	}
}

func TestCanEdit(t *testing.T) {
	user := &User{ID: 1}
	photo := &Photo{ID: 1, OwnerID: 1}

	if photo.CanEdit(user) {
		t.Error("Non-authenticated should not be able to edit")
	}

	user.IsAuthenticated = true

	if !photo.CanEdit(user) {
		t.Error("User should be able to edit")
	}

	photo.OwnerID = 2

	if photo.CanEdit(user) {
		t.Error("User should not be able to edit")
	}

	user.IsAdmin = true
	if !photo.CanEdit(user) {
		t.Error("Admin should be able to edit")
	}
}

func TestHasVoted(t *testing.T) {

	u := &User{}
	if u.HasVoted(1) {
		t.Error("The user has not voted yet")
	}

	u.RegisterVote(1)
	if !u.HasVoted(1) {
		t.Error("The user should have voted")
	}
}
