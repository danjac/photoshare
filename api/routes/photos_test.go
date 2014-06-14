package routes

import (
	"github.com/danjac/photoshare/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MakeMockAppContext(user *models.User) *AppContext {

	req := &http.Request{}
	params := make(map[string]string)
	res := httptest.NewRecorder()

	return &AppContext{req, res, params, user}

}

type MockPhotoManager struct {
}

func (m *MockPhotoManager) Get(photoID string) (*models.Photo, error) {
	return nil, nil
}

func (m *MockPhotoManager) GetDetail(photoID string) (*models.PhotoDetail, error) {
	return nil, nil
}

func (m *MockPhotoManager) All(pageNum int64) ([]models.Photo, error) {
	return []models.Photo{}, nil
}

func (m *MockPhotoManager) ByOwnerID(pageNum int64, ownerID string) ([]models.Photo, error) {
	return []models.Photo{}, nil
}

func (m *MockPhotoManager) Search(pageNum int64, q string) ([]models.Photo, error) {
	return []models.Photo{}, nil
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

func TestGetPhotos(t *testing.T) {
	photoMgr = &MockPhotoManager{}
	c := MakeMockAppContext(nil)
	if err := getPhotos(c); err != nil {
		t.Error(err)
	}

}
