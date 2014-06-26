package routes

import (
	"encoding/json"
	"github.com/danjac/photoshare/api/models"
	"github.com/zenazn/goji/web"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAnonymousSession struct {}

func (m *MockAnonymousSession) GetCurrentUser(r *http.Request) (*models.User, error) {
	return &models.User{}, nil
}

func (m *MockAnonymousSession) Login(w http.ResponseWriter, user *models.User) (string, error) {
	return "", nil
}

func (m *MockAnonymousSession) Login(w http.ResponseWriter) (string, error) {
	return "", nil
)

type MockPhotoManager struct {
}

func (m *MockPhotoManager) Get(photoID string) (*models.Photo, error) {
	return nil, nil
}

func (m *MockPhotoManager) GetDetail(photoID string, user *models.User) (*models.PhotoDetail, error) {
	return nil, nil
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

func TestGetPhotoDetail(t *testing T) {

	req := &http.Request{}
	res := httptest.NewRecorder()

	sessionMgr = &MockAnonymousSessionManager{}
	photoMgr = &MockPhotoManager{}

	getPhotoDetail(web.C{}, res, req)
	value := &models.PhotoDetail{}
	parseJsonBody(res, value)
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
