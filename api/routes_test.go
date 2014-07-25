package api

import (
	"database/sql"
	"encoding/json"
	"github.com/zenazn/goji/web"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockCache struct{}

func (m *mockCache) Set(key string, obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (m *mockCache) DeleteAll() error {
	return nil
}

func (m *mockCache) Get(key string, fn func() (interface{}, error)) (interface{}, error) {
	return fn()
}

func (m *mockCache) Render(w http.ResponseWriter, status int, key string, fn func() (interface{}, error)) error {
	obj, err := fn()
	if err != nil {
		return err
	}
	value, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	return writeBody(w, value, status, "application/json")
}

type mockSessionManager struct {
}

func (m *mockSessionManager) ReadToken(r *http.Request) (int64, error) {
	return 0, nil
}

func (m *mockSessionManager) WriteToken(w http.ResponseWriter, userID int64) error {
	return nil
}

type mockPhotoDataStore struct {
}

func (m *mockPhotoDataStore) Get(photoID int64) (*Photo, error) {
	return nil, sql.ErrNoRows
}

func (m *mockPhotoDataStore) GetDetail(photoID int64, user *User) (*PhotoDetail, error) {
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
	return photo, nil
}

func (m *mockPhotoDataStore) All(page *Page, orderBy string) (*PhotoList, error) {
	item := &Photo{
		ID:      1,
		Title:   "test",
		OwnerID: 1,
	}
	photos := []Photo{*item}
	return NewPhotoList(photos, 1, 1), nil
}

func (m *mockPhotoDataStore) ByOwnerID(page *Page, ownerID int64) (*PhotoList, error) {
	return &PhotoList{}, nil
}

func (m *mockPhotoDataStore) Search(page *Page, q string) (*PhotoList, error) {
	return &PhotoList{}, nil
}

func (m *mockPhotoDataStore) UpdateTags(photo *Photo) error {
	return nil
}

func (m *mockPhotoDataStore) GetTagCounts() ([]TagCount, error) {
	return []TagCount{}, nil
}

func (m *mockPhotoDataStore) Delete(photo *Photo) error {
	return nil
}

func (m *mockPhotoDataStore) Insert(photo *Photo) error {
	return nil
}

func (m *mockPhotoDataStore) Update(photo *Photo) error {
	return nil
}

type emptyPhotoDataStore struct {
	mockPhotoDataStore
}

func (m *emptyPhotoDataStore) All(page *Page, orderBy string) (*PhotoList, error) {
	var photos []Photo
	return &PhotoList{photos, 0, 1, 0}, nil
}

func (m *emptyPhotoDataStore) GetDetail(photoID int64, user *User) (*PhotoDetail, error) {
	return nil, sql.ErrNoRows
}

// should return a 404
func TestGetPhotoDetailIfNone(t *testing.T) {
	req := &http.Request{}
	res := httptest.NewRecorder()
	c := web.C{}
	c.Env = make(map[string]interface{})

	a := &AppContext{
		sessionMgr: &mockSessionManager{},
		ds:         &DataStores{photos: &emptyPhotoDataStore{}},
	}

	err := a.photoDetail(c, res, req)
	if err != sql.ErrNoRows {
		t.Fail()
	}
}

func TestGetPhotoDetail(t *testing.T) {

	req, _ := http.NewRequest("GET", "http://localhost/api/photos/1", nil)
	res := httptest.NewRecorder()
	c := web.C{}
	c.Env = make(map[string]interface{})
	c.URLParams = make(map[string]string)
	c.URLParams["id"] = "1"

	a := &AppContext{
		sessionMgr: &mockSessionManager{},
		ds:         &DataStores{photos: &mockPhotoDataStore{}},
	}

	c.Env["user"] = &User{}

	a.photoDetail(c, res, req)
	value := &PhotoDetail{}
	parseJSONBody(res, value)
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

	a := &AppContext{
		ds:    &DataStores{photos: &mockPhotoDataStore{}},
		cache: &mockCache{},
	}

	a.getPhotos(web.C{}, res, req)
	value := &PhotoList{}
	parseJSONBody(res, value)
	if value.Total != 1 {
		t.Fail()
	}

}
