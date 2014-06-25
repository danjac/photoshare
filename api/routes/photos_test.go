package routes

import (
	"fmt"
	"github.com/danjac/photoshare/api/models"
	"github.com/zenazn/goji/web"
	"net/http"
	"net/http/httptest"
	"testing"
)

func MakeMockContext(user *models.User) (web.C, http.ResponseWriter, *http.Request) {

	req := &http.Request{}
	c := web.C{}
	c.URLParams = make(map[string]string)
	res := httptest.NewRecorder()

	return c, res, req

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
	c, w, r := MakeMockContext(nil)
	getPhotos(c, w, r)
	fmt.Println(w)

}
