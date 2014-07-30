package photoshare

import (
	"database/sql"
	"testing"
)

func TestGetIfNotNone(t *testing.T) {

	cfg, _ := newConfigurator()
	tdb := makeTestDB(cfg)
	defer tdb.clean()

	datamapper, _ := newDataMapper(tdb.dbMap.Db, false)

	user := &user{Name: "tester", Email: "tester@gmail.com", Password: "test"}

	if err := datamapper.createUser(user); err != nil {
		t.Error(err)
		return
	}
	photo := &photo{Title: "test", OwnerID: user.ID, Filename: "test.jpg"}
	if err := datamapper.createPhoto(photo); err != nil {
		t.Error(err)
		return
	}

	photo, err := datamapper.getPhoto(photo.ID)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestGetIfNone(t *testing.T) {

	cfg, _ := newConfigurator()
	tdb := makeTestDB(cfg)
	defer tdb.clean()

	datamapper, _ := newDataMapper(tdb.dbMap.Db, false)

	_, err := datamapper.getPhoto(1)
	if err != sql.ErrNoRows {
		t.Error(err)
		return
	}

}

func TestSearchPhotos(t *testing.T) {
	cfg, _ := newConfigurator()
	tdb := makeTestDB(cfg)
	defer tdb.clean()

	datamapper, _ := newDataMapper(tdb.dbMap.Db, false)

	user := &user{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := datamapper.createUser(user); err != nil {
		t.Error(err)
		return
	}
	photo := &photo{Title: "test", OwnerID: user.ID, Filename: "test.jpg"}
	if err := datamapper.createPhoto(photo); err != nil {
		t.Error(err)
		return
	}

	result, err := datamapper.searchPhotos(newPage(1), "test")
	if err != nil {
		t.Error(err)
		return
	}

	if len(result.Items) != 1 {
		t.Error("There should be 1 photo")
	}
}
func TestAllPhotos(t *testing.T) {
	cfg, _ := newConfigurator()
	tdb := makeTestDB(cfg)
	defer tdb.clean()

	datamapper, _ := newDataMapper(tdb.dbMap.Db, false)

	user := &user{Name: "tester", Email: "tester@gmail.com", Password: "test"}
	if err := datamapper.createUser(user); err != nil {
		t.Error(err)
		return
	}
	photo := &photo{Title: "test", OwnerID: user.ID, Filename: "test.jpg"}
	if err := datamapper.createPhoto(photo); err != nil {
		t.Error(err)
		return
	}

	result, err := datamapper.getPhotos(newPage(1), "")
	if err != nil {
		t.Error(err)
		return
	}

	if len(result.Items) != 1 {
		t.Error("There should be 1 photo")
	}
}

func TestCanEdit(t *testing.T) {
	user := &user{ID: 1}
	photo := &photo{ID: 1, OwnerID: 1}

	if photo.canEdit(user) {
		t.Error("Non-authenticated should not be able to edit")
	}

	user.IsAuthenticated = true

	if !photo.canEdit(user) {
		t.Error("User should be able to edit")
	}

	photo.OwnerID = 2

	if photo.canEdit(user) {
		t.Error("User should not be able to edit")
	}

	user.IsAdmin = true
	if !photo.canEdit(user) {
		t.Error("Admin should be able to edit")
	}
}

func TestHasVoted(t *testing.T) {

	u := &user{}
	if u.hasVoted(1) {
		t.Error("The user has not voted yet")
	}

	u.registerVote(1)
	if !u.hasVoted(1) {
		t.Error("The user should have voted")
	}
}
