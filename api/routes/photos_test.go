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

func TestGetPhotos(t *testing.T) {

	//c := MakeMockAppContext(nil)
	//getPhotos(c)

}
