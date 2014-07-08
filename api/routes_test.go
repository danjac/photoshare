package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPhotoManager struct {
}

func (m *mockPhotoManager) Get(photoID int64) (*Photo, bool, error) {
	return nil, false, nil
}

func (m *mockPhotoManager) GetDetail(photoID int64, user *User) (*PhotoDetail, bool, error) {
	canEdit := user.ID == 1
	photo := &PhotoDetail{
		Photo: Photo{
			ID:      1,
			Title:   "test",
			OwnerID: 1,
		},
		OwnerName: "tester",
		Permissions: &Permissions{
			Edit: canEdit,
		},
	}
	return photo, true, nil
}

func (m *mockPhotoManager) All(page *Page, orderBy string) (*PhotoList, error) {
	item := &Photo{
		ID:      1,
		Title:   "test",
		OwnerID: 1,
	}
	photos := []Photo{*item}
	return NewPhotoList(photos, 1, 1), nil
}

func (m *mockPhotoManager) ByOwnerID(page *Page, ownerID int64) (*PhotoList, error) {
	return &PhotoList{}, nil
}

func (m *mockPhotoManager) Search(page *Page, q string) (*PhotoList, error) {
	return &PhotoList{}, nil
}

func (m *mockPhotoManager) UpdateTags(photo *Photo) error {
	return nil
}

func (m *mockPhotoManager) GetTagCounts() ([]TagCount, error) {
	return []TagCount{}, nil
}

func (m *mockPhotoManager) Delete(photo *Photo) error {
	return nil
}

func (m *mockPhotoManager) Insert(photo *Photo) error {
	return nil
}

func (m *mockPhotoManager) Update(photo *Photo) error {
	return nil
}

type emptyPhotoManager struct {
	mockPhotoManager
}

func (m *emptyPhotoManager) All(page *Page, orderBy string) (*PhotoList, error) {
	var photos []Photo
	return &PhotoList{photos, 0, 1, 0}, nil
}

func (m *emptyPhotoManager) GetDetail(photoID int64, user *User) (*PhotoDetail, bool, error) {
	return nil, false, nil
}

// should return a 404
func TestGetPhotoDetailIfNone(t *testing.T) {
	req := &http.Request{}
	res := httptest.NewRecorder()

	getCurrentUser = func(r *http.Request) (*User, error) {
		return &User{}, nil
	}

	photoMgr = &emptyPhotoManager{}

	photoDetail(res, req)
	if res.Code != 404 {
		t.Fail()
	}
}

func TestGetPhotoDetailWithBadID(t *testing.T) {
	req := &http.Request{}
	res := httptest.NewRecorder()

	getCurrentUser = func(r *http.Request) (*User, error) {
		return &User{}, nil
	}

	photoMgr = &mockPhotoManager{}
	photoDetail(res, req)
	if res.Code != 404 {
		t.Fatal("Should be a 404")
	}

}

func TestGetPhotoDetail(t *testing.T) {

	req, _ := http.NewRequest("GET", "http://localhost/api/photos/1", nil)
	res := httptest.NewRecorder()

	getCurrentUser = func(r *http.Request) (*User, error) {
		return &User{}, nil
	}

	routeParamInt64 = func(r *http.Request, name string) (int64, error) {
		return int64(1), nil
	}

	photoMgr = &mockPhotoManager{}

	photoDetail(res, req)
	value := &PhotoDetail{}
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
	getPhotos(res, req)
	value := &PhotoList{}
	parseJsonBody(res, value)
	if value.Total != 1 {
		t.Fail()
	}

}
