package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) error {

	if err := session.Logout(w); err != nil {
		return err
	}

	return render(w, http.StatusOK, "Logged out")

}

// return current logged in user, or 401
func authenticate(w http.ResponseWriter, r *http.Request) error {

	user, err := session.GetCurrentUser(r)
	if err != nil {
		return err
	}
	var status int
	if user == nil {
		status = http.StatusNotFound
	} else {
		status = http.StatusOK
	}

	return render(w, status, user)
}

func login(w http.ResponseWriter, r *http.Request) error {

	auth := &models.Authenticator{}
	if err := parseJSON(r, auth); err != nil {
		return err
	}

	user, err := auth.Identify()
	if err != nil {
		if err == models.MissingLoginFields {
			return render(w, http.StatusBadRequest, "Missing email or password")
		}
		return err
	}

	if user == nil {
		return render(w, http.StatusBadRequest, "Invalid email or password")
	}

	if err := session.Login(w, user); err != nil {
		return err
	}

	return render(w, http.StatusOK, user)
}
