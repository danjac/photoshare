package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"net/http"
)

func logout(w http.ResponseWriter, r *http.Request) {

	if err := session.Logout(w); err != nil {
		panic(err)
	}

	render(w, http.StatusOK, "Logged out")

}

// return current logged in user, or 401
func authenticate(w http.ResponseWriter, r *http.Request) {

	user, err := session.GetCurrentUser(r)
	if err != nil {
		panic(err)
	}
	var status int
	if user == nil {
		status = http.StatusNotFound
	} else {
		status = http.StatusOK
	}

	render(w, status, user)
}

func login(w http.ResponseWriter, r *http.Request) {

	auth := &models.Authenticator{}
	if err := parseJSON(r, auth); err != nil {
		panic(err)
	}

	user, err := auth.Identify()
	if err != nil {
		if err == models.MissingLoginFields {
			render(w, http.StatusBadRequest, "Missing email or password")
			return
		}
		panic(err)
	}

	if user == nil {
		render(w, http.StatusBadRequest, "Invalid email or password")
		return
	}

	if err := session.Login(w, user); err != nil {
		panic(err)
	}

	render(w, http.StatusOK, user)
}
