package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"net/http"
)

func signup(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}

	if err := parseJSON(r, user); err != nil {
		panic(err)
	}

	if result, err := user.Validate(); err != nil || !result.OK {
		if err != nil {
			panic(err)
		}
		render(w, http.StatusBadRequest, result)
	}

	if err := user.Insert(); err != nil {
		panic(err)
	}

	if err := session.Login(w, user); err != nil {
		panic(err)
	}

	render(w, http.StatusOK, user)

}
