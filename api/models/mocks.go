package models

type MockPhotoManager struct{}

func (mgr *MockPhotoManager) Create(*Photo) error {
    return nil
}

func (mgr *MockPhotoManager) Update(*Photo) error {
    return nil
}

func (mgr *MockPhotoManager) Delete(*Photo) error {
    return nil
}

func (mgr *MockPhotoManager) Get(photoID string) (*Photo, error) {
    return &Photo{ID: 1, Title: "test", Photo: "test.jpg"}, nil
}

func (mgr *MockPhotoManager) GetDetail(photoID string) (*PhotoDetail, error) {
    return &PhotoDetail{ID: 1, Title: "test", Photo: "test.jpg"}, nil
}


func (mgr *MockPhotoManager) GetAll(pageNum int64) ([]Photo, error) {
	return []Photo{
		Photo{ID: 1, Title: "test", Photo: "test.jpg"},
	}, nil
}


