package routes

import (
	"encoding/json"
	"github.com/danjac/photoshare/api/models"
	"github.com/zenazn/goji/web"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAnonymousSessionManager struct{}

func (m *MockAnonymousSessionManager) GetCurrentUser(r *http.Request) (*models.User, error) {
	return &models.User{}, nil
}

func (m *MockAnonymousSessionManager) Login(w http.ResponseWriter, user *models.User) (string, error) {
	return "", nil
}

func (m *MockAnonymousSessionManager) Logout(w http.ResponseWriter) (string, error) {
	return "", nil
}

type MockPhotoManager struct {
}

func (m *MockPhotoManager) Get(photoID string) (*models.Photo, error) {
	return nil, nil
}

func (m *MockPhotoManager) GetDetail(photoID string, user *models.User) (*models.PhotoDetail, error) {
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

func (m *MockPhotoManager) All(pageNum int64, orderBy string) (*models.PhotoList, error) {
	item := &models.Photo{
		ID:      1,
		Title:   "test",
		OwnerID: 1,
	}
	photos := []models.Photo{*item}
	return models.NewPhotoList(photos, 1, 1), nil
}

func (m *MockPhotoManager) ByOwnerID(pageNum int64, ownerID string) (*models.PhotoList, error) {
	return &models.PhotoList{}, nil
}

func (m *MockPhotoManager) Search(pageNum int64, q string) (*models.PhotoList, error) {
	return &models.PhotoList{}, nil
}

func (m *MockPhotoManager) UpdateTags(photo *models.Photo) error {
	return nil
}

func (m *MockPhotoManager) GetTagCounts() ([]models.TagCount, error) {
	return []models.TagCount{}, nil
}

func (m *MockPhotoManager) Delete(photo *models.Photo) error {
	return nil
}

func (m *MockPhotoManager) Insert(photo *models.Photo) error {
	return nil
}

func (m *MockPhotoManager) Update(photo *models.Photo) error {
	return nil
}

func parseJsonBody(res *httptest.ResponseRecorder, value interface{}) error {
	return json.Unmarshal([]byte(res.Body.String()), value)
}

func newContext() web.C {
	c := web.C{}
	c.Env = make(map[string]interface{})
	c.URLParams = make(map[string]string)
	return c
}

func TestGetPhotoDetail(t *testing.T) {

	req := &http.Request{}
	res := httptest.NewRecorder()
	c := newContext()

	sessionMgr = &MockAnonymousSessionManager{}
	photoMgr = &MockPhotoManager{}

	photoDetail(c, res, req)
	value := &models.PhotoDetail{}
	parseJsonBody(res, value)
	if value.Title != "test" {
		t.Fail()
	}
	if value.Permissions.Edit {
		t.Fail()
	}
}

func TestGetPhotos(t *testing.T) {

	req := &http.Request{}
	res := httptest.NewRecorder()

	photoMgr = &MockPhotoManager{}
	getPhotos(web.C{}, res, req)
	value := &models.PhotoList{}
	parseJsonBody(res, value)
	if value.Total != 1 {
		t.Fail()
	}

}
