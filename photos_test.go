package photoshare

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func (m *mockSessionManager) createToken(userID int64) (string, error) {
	return strconv.FormatInt(userID, 10), nil
}

func (m *mockSessionManager) writeToken(w http.ResponseWriter, userID int64) error {
	return nil
}

type mockDataMapper struct {
}

func (m *mockDataMapper) getPhoto(photoID int64) (*photo, error) {
	return nil, sql.ErrNoRows
}

func (m *mockDataMapper) getPhotoDetail(photoID int64, user *user) (*photoDetail, error) {
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

func (m *mockDataMapper) getPhotos(page *page, orderBy string) (*photoList, error) {
	item := &photo{
		ID:      1,
		Title:   "test",
		OwnerID: 1,
	}
	photos := []photo{*item}
	return newPhotoList(photos, 1, 1), nil
}

func (m *mockDataMapper) getPhotosByOwnerID(page *page, ownerID int64) (*photoList, error) {
	return &photoList{}, nil
}

func (m *mockDataMapper) searchPhotos(page *page, q string) (*photoList, error) {
	return &photoList{}, nil
}

func (m *mockDataMapper) getTagCounts() ([]tagCount, error) {
	return []tagCount{}, nil
}

func (m *mockDataMapper) getActiveUser(userID int64) (*user, error) {
	return &user{}, nil
}

func (m *mockDataMapper) getUserByEmail(email string) (*user, error) {
	return &user{}, nil
}

func (m *mockDataMapper) isUserNameAvailable(user *user) (bool, error) {
	return true, nil
}

func (m *mockDataMapper) isUserEmailAvailable(user *user) (bool, error) {
	return true, nil
}

func (m *mockDataMapper) getUserByNameOrEmail(identifier string) (*user, error) {
	return &user{}, nil
}

func (m *mockDataMapper) getUserByRecoveryCode(code string) (*user, error) {
	return &user{}, nil
}

func (m *mockDataMapper) createPhoto(_ *photo) error {
	return nil
}

func (m *mockDataMapper) removePhoto(_ *photo) error {
	return nil
}

func (m *mockDataMapper) updatePhoto(_ *photo) error {
	return nil
}

func (m *mockDataMapper) updateTags(_ *photo) error {
	return nil
}

func (m *mockDataMapper) createUser(_ *user) error {
	return nil
}

func (m *mockDataMapper) updateUser(_ *user) error {
	return nil
}

func (m *mockDataMapper) updateMany(items ...interface{}) error {
	return nil
}

type emptyDataStore struct {
	mockDataMapper
}

func (m *emptyDataStore) getPhotos(page *page, orderBy string) (*photoList, error) {
	var photos []photo
	return &photoList{photos, 0, 1, 0}, nil
}

func (m *emptyDataStore) getPhotoDetail(photoID int64, user *user) (*photoDetail, error) {
	return nil, sql.ErrNoRows
}

// should return a 404
func TestGetPhotoDetailIfNone(t *testing.T) {
	req := &http.Request{}
	res := httptest.NewRecorder()

	cfg := &configurator{
		session:    &mockSessionManager{},
		datamapper: &emptyDataStore{},
	}

	c := &context{
		configurator: cfg,
		params:       &params{make(map[string]string)},
	}

	err := getPhotoDetail(c, res, req)
	if err != sql.ErrNoRows {
		t.Fail()
	}
}

func TestGetPhotoDetail(t *testing.T) {

	req, _ := http.NewRequest("GET", "http://localhost/api/photos/1", nil)
	res := httptest.NewRecorder()
	p := &params{make(map[string]string)}
	p.vars["id"] = "1"

	cfg := &configurator{
		session:    &mockSessionManager{},
		datamapper: &mockDataMapper{},
	}

	c := &context{
		configurator: cfg,
		params:       p,
	}

	getPhotoDetail(c, res, req)
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

	cfg := &configurator{
		datamapper: &mockDataMapper{},
		cache:      &mockCache{},
	}

	c := &context{
		configurator: cfg,
		params:       &params{},
	}

	getPhotos(c, res, req)
	value := &photoList{}
	parseJSONBody(res, value)
	if value.Total != 1 {
		t.Fail()
	}

}
