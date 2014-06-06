package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/render"
	"github.com/danjac/photoshare/api/session"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) error {

	if err := session.Logout(w); err != nil {
		return err
	}

	return render.Status(w, http.StatusOK, "Logged out")

}

// return current logged in user, or 401
func authenticate(w http.ResponseWriter, r *http.Request) error {

	user, err := session.GetCurrentUser(r)
	if err != nil {
		return err
	}

	var status int

	if user != nil {
		status = http.StatusOK
	} else {
		status = http.StatusNotFound
	}

	return render.JSON(w, status, user)
}

func login(w http.ResponseWriter, r *http.Request) error {

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		return render.Status(w, http.StatusBadRequest, "Email or password missing")
	}

	user, err := models.Authenticate(email, password)
	if err != nil {
		return err
	}

	if user == nil {
		return render.Status(w, http.StatusBadRequest, "Invalid email or password")
	}

	if err := session.Login(w, user); err != nil {
		return err
	}

	return render.JSON(w, http.StatusOK, user)
}
