package routes

import (
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/validation"
	"github.com/zenazn/goji/web"
	"net/http"
	"strings"
)

var getUserValidator = func(user *models.User) validation.Validator {
	return validation.NewUserValidator(user)
}

var getCurrentUser = func(c web.C, r *http.Request) (*models.User, error) {

	obj, ok := c.Env["user"]
	if ok {
		return obj.(*models.User), nil
	}

	user, err := sessionMgr.GetCurrentUser(r)
	if err != nil {
		return nil, err
	}

	c.Env["user"] = user
	return user, nil
}

// Basic user session info
type SessionInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"isAdmin"`
	LoggedIn bool   `json:"loggedIn"`
}

func newSessionInfo(user *models.User) *SessionInfo {
	if user == nil || user.ID == 0 || !user.IsAuthenticated {
		return &SessionInfo{}
	}

	return &SessionInfo{user.ID, user.Name, user.IsAdmin, true}
}

func logout(c web.C, w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(c, r)
	if !user.IsAuthenticated {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if _, err = sessionMgr.Logout(w); err != nil {
		panic(err)
	}

	sendMessage(&Message{user.Name, "", 0, "logout"})
	writeJSON(w, newSessionInfo(&models.User{}), http.StatusOK)

}

func authenticate(c web.C, w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(c, r)
	if err != nil {
		panic(err)
	}

	writeJSON(w, newSessionInfo(user), http.StatusOK)
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
		writeString(w, "Missing login details", http.StatusBadRequest)
		return
	}

	user, err := userMgr.Authenticate(s.Identifier, s.Password)

	if err != nil {
		panic(err)
	}
	if !user.IsAuthenticated {
		writeString(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	if _, err := sessionMgr.Login(w, user); err != nil {
		panic(err)
	}
	sendMessage(&Message{user.Name, "", 0, "login"})
	writeJSON(w, newSessionInfo(user), http.StatusOK)
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

	validator := getUserValidator(user)

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			panic(err)
		}
		writeJSON(w, result, http.StatusBadRequest)
		return
	}

	if err := userMgr.Insert(user); err != nil {
		panic(err)
	}

	if _, err := sessionMgr.Login(w, user); err != nil {
		panic(err)
	}

	user.IsAuthenticated = true

	writeJSON(w, newSessionInfo(user), http.StatusOK)

}
