package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/render"
	"github.com/danjac/photoshare/api/session"
	"net/http"
)

func signup(w http.ResponseWriter, r *http.Request) {

	user := models.NewUser(
		r.FormValue("name"),
		r.FormValue("email"),
		r.FormValue("password"),
	)

	if result, err := user.Validate(); err != nil || !result.OK {
		if err != nil {
			render.Error(w, r, err)
			return
		}
		render.JSON(w, http.StatusBadRequest, result)
		return
	}

	if err := user.Save(); err != nil {
		render.Error(w, r, err)
		return
	}

	if err := session.Login(w, user); err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, http.StatusOK, user)

}
