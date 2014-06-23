package routes

import (
	"github.com/danjac/photoshare/api/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MakeMockContext(user *models.User) *Context {

	req := &http.Request{}
	params := make(map[string]string)
	res := httptest.NewRecorder()

	return &Context{req, res, params, user, nil}

}

type MockPhotoManager struct {
}

func (m *MockPhotoManager) Get(photoID string) (*models.Photo, error) {
	return nil, nil
}

func (m *MockPhotoManager) GetDetail(photoID string, user *models.User) (*models.PhotoDetail, error) {
	return nil, nil
}

func (m *MockPhotoManager) All(pageNum int64, orderBy string) (*models.PhotoList, error) {
	photos := make([]models.Photo, 0, 0)
	return models.NewPhotoList(photos, 0, 0), nil
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

func TestGetPhotos(t *testing.T) {
	photoMgr = &MockPhotoManager{}
	c := MakeMockContext(nil)
	result := getPhotos(c)
	if result.Error != nil {
		t.Error(result.Error)
	}

}
