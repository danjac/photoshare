package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/zenazn/goji/web"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPhotoManager struct {
}

func (m *mockPhotoManager) Get(photoID int64) (*models.Photo, error) {
	return nil, nil
}

func (m *mockPhotoManager) GetDetail(photoID int64, user *models.User) (*models.PhotoDetail, error) {
	canEdit := user.ID == 1
	photo := &models.PhotoDetail{
		Photo: models.Photo{
			ID:      1,
			Title:   "test",
			OwnerID: 1,
		},
		OwnerName: "tester",
		Permissions: &models.Permissions{
			Edit: canEdit,
		},
	}
	return photo, nil
}

func (m *mockPhotoManager) All(pageNum int64, orderBy string) (*models.PhotoList, error) {
	item := &models.Photo{
		ID:      1,
		Title:   "test",
		OwnerID: 1,
	}
	photos := []models.Photo{*item}
	return models.NewPhotoList(photos, 1, 1), nil
}

func (m *mockPhotoManager) ByOwnerID(pageNum int64, ownerID int64) (*models.PhotoList, error) {
	return &models.PhotoList{}, nil
}

func (m *mockPhotoManager) Search(pageNum int64, q string) (*models.PhotoList, error) {
	return &models.PhotoList{}, nil
}

func (m *mockPhotoManager) UpdateTags(photo *models.Photo) error {
	return nil
}

func (m *mockPhotoManager) GetTagCounts() ([]models.TagCount, error) {
	return []models.TagCount{}, nil
}

func (m *mockPhotoManager) Delete(photo *models.Photo) error {
	return nil
}

func (m *mockPhotoManager) Insert(photo *models.Photo) error {
	return nil
}

func (m *mockPhotoManager) Update(photo *models.Photo) error {
	return nil
}

type emptyPhotoManager struct {
	mockPhotoManager
}

func (m *emptyPhotoManager) All(pageNum int64, orderBy string) (*models.PhotoList, error) {
	var photos []models.Photo
	return &models.PhotoList{photos, 0, 1, 0}, nil
}

func (m *emptyPhotoManager) GetDetail(photoID int64, user *models.User) (*models.PhotoDetail, error) {
	return nil, nil
}

// should return a 404
func TestGetPhotoDetailIfNone(t *testing.T) {
	req := &http.Request{}
	res := httptest.NewRecorder()
	c := newContext()

	getCurrentUser = func(c web.C, r *http.Request) (*models.User, error) {
		return &models.User{}, nil
	}

	photoMgr = &emptyPhotoManager{}

	photoDetail(c, res, req)
	if res.Code != 404 {
		t.Fail()
	}
}

func TestGetPhotoDetailWithBadID(t *testing.T) {
	req := &http.Request{}
	res := httptest.NewRecorder()
	c := newContext()
	c.URLParams["id"] = "fiddlesticks"

	getCurrentUser = func(c web.C, r *http.Request) (*models.User, error) {
		return &models.User{}, nil
	}

	photoMgr = &mockPhotoManager{}
	photoDetail(c, res, req)
	if res.Code != 404 {
		t.Fatal("Should be a 404")
	}

}

func TestGetPhotoDetail(t *testing.T) {

	req := &http.Request{}
	res := httptest.NewRecorder()
	c := newContext()
	c.URLParams["id"] = "1"

	getCurrentUser = func(c web.C, r *http.Request) (*models.User, error) {
		return &models.User{}, nil
	}

	photoMgr = &mockPhotoManager{}

	photoDetail(c, res, req)
	value := &models.PhotoDetail{}
	parseJsonBody(res, value)
	if res.Code != 200 {
		t.Fatal("Photo not found")
	}
	if value.Title != "test" {
		t.Fatal("Title should be test")
	}
	if value.Permissions.Edit {
		t.Fatal("User should have edit permission")
	}
}

func TestGetPhotos(t *testing.T) {

	req := &http.Request{}
	res := httptest.NewRecorder()

	photoMgr = &mockPhotoManager{}
	getPhotos(web.C{}, res, req)
	value := &models.PhotoList{}
	parseJsonBody(res, value)
	if value.Total != 1 {
		t.Fail()
	}

}
