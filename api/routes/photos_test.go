package routes

import (
	"github.com/danjac/photoshare/api/models"
	"testing"
    "net/http"
    "net/http/httptest"
)

func MakeMockAppContext(user *models.User) *AppContext {

    req := &http.Request{}
    params := make(map[string]string)
    res := httptest.NewRecorder()
    var photoMgr = &models.MockPhotoManager{}

    return &AppContext{req, res, params, user, photoMgr}
    
}

func TestGetPhotos(t *testing.T) {

    c := MakeMockAppContext(nil)
    getPhotos(c)

}
