package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/render"
	"github.com/danjac/photoshare/api/session"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {

	if err := session.Logout(w); err != nil {
		render.Error(w, r, err)
		return
	}

	render.Status(w, http.StatusOK, "Logged out")

}

// return current logged in user, or 401
func authenticate(w http.ResponseWriter, r *http.Request) {

	user, err := session.GetCurrentUser(r)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	var status int

	if user != nil {
		status = http.StatusOK
	} else {
		status = http.StatusNotFound
	}

	render.JSON(w, status, user)
}

func login(w http.ResponseWriter, r *http.Request) {

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		render.Status(w, http.StatusBadRequest, "Email or password missing")
		return
	}

	user, err := models.Authenticate(email, password)
	if err != nil {
		render.Error(w, r, err)
		return
	}

	if user == nil {
		render.Status(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	if err := session.Login(w, user); err != nil {
		render.Error(w, r, err)
		return
	}

	render.JSON(w, http.StatusOK, user)
}
