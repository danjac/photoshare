package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/render"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/validation"
	"github.com/zenazn/goji/web"
	"net/http"
	"strings"
)

func logout(c web.C, w http.ResponseWriter, r *http.Request) {

	if _, err := session.Logout(w); err != nil {
		panic(err)
	}

	render.JSON(w, session.NewSessionInfo(&models.User{}), http.StatusOK)

}

func authenticate(c web.C, w http.ResponseWriter, r *http.Request) {

	user, err := session.GetCurrentUser(c, r)
	if err != nil {
		panic(err)
	}

	render.JSON(w, session.NewSessionInfo(user), http.StatusOK)
}

func login(c web.C, w http.ResponseWriter, r *http.Request) {

	s := &struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{}

	if err := parseJSON(r, s); err != nil {
		panic(err)
	}

	if s.Identifier == "" || s.Password == "" {
		render.String(w, "Missing login details", http.StatusBadRequest)
		return
	}

	user, err := userMgr.Authenticate(s.Identifier, s.Password)

	if err != nil {
		panic(err)
	}
	if !user.IsAuthenticated {
		render.String(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	if _, err := session.Login(w, user); err != nil {
		panic(err)
	}
	render.JSON(w, session.NewSessionInfo(user), http.StatusOK)
}

func signup(c web.C, w http.ResponseWriter, r *http.Request) {

	user := &models.User{}

	if err := parseJSON(r, user); err != nil {
		panic(err)
	}

	// ensure nobody tries to make themselves an admin
	user.IsAdmin = false

	// email should always be lower case
	user.Email = strings.ToLower(user.Email)

	validator := validation.NewUserValidator(user)

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			panic(err)
		}
		render.JSON(w, result, http.StatusBadRequest)
		return
	}

	if err := userMgr.Insert(user); err != nil {
		panic(err)
	}

	if _, err := session.Login(w, user); err != nil {
		panic(err)
	}

	user.IsAuthenticated = true

	render.JSON(w, session.NewSessionInfo(user), http.StatusOK)

}
