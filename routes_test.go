package photoshare

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockCache struct{}

func (m *mockCache) set(key string, obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (m *mockCache) clear() error {
	return nil
}

func (m *mockCache) get(key string, fn func() (interface{}, error)) (interface{}, error) {
	return fn()
}

func (m *mockCache) render(w http.ResponseWriter, status int, key string, fn func() (interface{}, error)) error {
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

func (m *mockSessionManager) readToken(r *http.Request) (int64, error) {
	return 0, nil
}

func (m *mockSessionManager) writeToken(w http.ResponseWriter, userID int64) error {
	return nil
}

type mockPhotoDataManager struct {
}

func (m *mockPhotoDataManager) get(photoID int64) (*photo, error) {
	return nil, sql.ErrNoRows
}

func (m *mockPhotoDataManager) getDetail(photoID int64, user *user) (*photoDetail, error) {
	canEdit := user.ID == 1
	photo := &photoDetail{
		photo: photo{
			ID:      1,
			Title:   "test",
			OwnerID: 1,
		},
		OwnerName: "tester",
		Permissions: &permissions{
			Edit: canEdit,
		},
	}
	return photo, nil
}

func (m *mockPhotoDataManager) all(page *page, orderBy string) (*photoList, error) {
	item := &photo{
		ID:      1,
		Title:   "test",
		OwnerID: 1,
	}
	photos := []photo{*item}
	return newPhotoList(photos, 1, 1), nil
}

func (m *mockPhotoDataManager) byOwnerID(page *page, ownerID int64) (*photoList, error) {
	return &photoList{}, nil
}

func (m *mockPhotoDataManager) search(page *page, q string) (*photoList, error) {
	return &photoList{}, nil
}

func (m *mockPhotoDataManager) updateTags(photo *photo) error {
	return nil
}

func (m *mockPhotoDataManager) updateVotes(photo *photo, user *user) error {
	return nil
}

func (m *mockPhotoDataManager) getTagCounts() ([]tagCount, error) {
	return []tagCount{}, nil
}

func (m *mockPhotoDataManager) remove(photo *photo) error {
	return nil
}

func (m *mockPhotoDataManager) create(photo *photo) error {
	return nil
}

func (m *mockPhotoDataManager) update(photo *photo) error {
	return nil
}

type emptyPhotoDataManager struct {
	mockPhotoDataManager
}

func (m *emptyPhotoDataManager) all(page *page, orderBy string) (*photoList, error) {
	var photos []photo
	return &photoList{photos, 0, 1, 0}, nil
}

func (m *emptyPhotoDataManager) getDetail(photoID int64, user *user) (*photoDetail, error) {
	return nil, sql.ErrNoRows
}

// should return a 404
func TestGetPhotoDetailIfNone(t *testing.T) {
	req := &http.Request{}
	res := httptest.NewRecorder()
	p := &params{make(map[string]string)}

	a := &appContext{
		session:   &mockSessionManager{},
		datastore: &dataStore{photos: &emptyPhotoDataManager{}},
	}

	err := getPhotoDetail(a, res, req, p)
	if err != sql.ErrNoRows {
		t.Fail()
	}
}

func TestGetPhotoDetail(t *testing.T) {

	req, _ := http.NewRequest("GET", "http://localhost/api/photos/1", nil)
	res := httptest.NewRecorder()
	p := &params{make(map[string]string)}
	p.vars["id"] = "1"

	a := &appContext{
		session:   &mockSessionManager{},
		datastore: &dataStore{photos: &mockPhotoDataManager{}},
	}

	getPhotoDetail(a, res, req, p)
	value := &photoDetail{}
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
	p := &params{}

	a := &appContext{
		datastore: &dataStore{photos: &mockPhotoDataManager{}},
		cache:     &mockCache{},
	}

	getPhotos(a, res, req, p)
	value := &photoList{}
	parseJSONBody(res, value)
	if value.Total != 1 {
		t.Fail()
	}

}
