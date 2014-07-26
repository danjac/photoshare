package photoshare

import (
	"net/http"
	"strings"
)

// Basic user session info
type sessionInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"isAdmin"`
	LoggedIn bool   `json:"loggedIn"`
}

func newSessionInfo(user *user) *sessionInfo {
	if user == nil || user.ID == 0 || !user.IsAuthenticated {
		return &sessionInfo{}
	}

	return &sessionInfo{user.ID, user.Name, user.IsAdmin, true}
}

func (a *appContext) authenticate(r *request, required bool) (*user, error) {

	var invalidLogin error

	if required {
		invalidLogin = &httpError{http.StatusUnauthorized, "You must be logged in"}
	}

	if r.user != nil {
		return r.user, nil
	}
	r.user = &user{}

	userID, err := a.sessionMgr.readToken(r)
	if err != nil {
		return r.user, err
	}
	if userID == 0 {
		return r.user, invalidLogin
	}
	r.user, err = a.ds.users.getActive(userID)
	if err != nil {
		if isErrSqlNoRows(err) {
			return r.user, invalidLogin
		}
		return r.user, err
	}
	r.user.IsAuthenticated = true

	return r.user, nil
}

func (a *appContext) logout(w http.ResponseWriter, r *request) error {

	u, err := a.authenticate(r, true)
	if err != nil {
		return err
	}

	if err := a.sessionMgr.writeToken(w, 0); err != nil {
		return err
	}

	sendMessage(&socketMessage{u.Name, "", 0, "logout"})
	return renderJSON(w, newSessionInfo(&user{}), http.StatusOK)

}

func (a *appContext) getSessionInfo(w http.ResponseWriter, r *request) error {

	user, err := a.authenticate(r, false)
	if err != nil {
		return err
	}

	return renderJSON(w, newSessionInfo(user), http.StatusOK)
}

func (a *appContext) login(w http.ResponseWriter, r *request) error {

	s := &struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{}

	var invalidLogin = &httpError{http.StatusBadRequest, "Invalid email or password"}

	if err := r.decodeJSON(s); err != nil {
		return err
	}

	if s.Identifier == "" || s.Password == "" {
		return invalidLogin
	}

	user, err := a.ds.users.getByNameOrEmail(s.Identifier)
	if err != nil {
		if isErrSqlNoRows(err) {
			return invalidLogin
		}
		return err
	}
	if !user.checkPassword(s.Password) {
		return invalidLogin
	}

	if err := a.sessionMgr.writeToken(w, user.ID); err != nil {
		return err
	}

	user.IsAuthenticated = true

	sendMessage(&socketMessage{user.Name, "", 0, "login"})
	return renderJSON(w, newSessionInfo(user), http.StatusCreated)
}

func (a *appContext) signup(w http.ResponseWriter, r *request) error {

	s := &struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := r.decodeJSON(s); err != nil {
		return err
	}

	user := &user{
		Name:     s.Name,
		Email:    strings.ToLower(s.Email),
		Password: s.Password,
	}

	if err := validate(newUserValidator(user, a.ds.users)); err != nil {
		return err
	}

	if err := a.ds.users.create(user); err != nil {
		return err
	}

	if err := a.sessionMgr.writeToken(w, user.ID); err != nil {
		return err
	}

	user.IsAuthenticated = true

	go func() {
		if err := a.mailer.sendWelcomeMail(user); err != nil {
			logError(err)
		}
	}()

	return renderJSON(w, newSessionInfo(user), http.StatusCreated)

}

func (a *appContext) changePassword(w http.ResponseWriter, r *request) error {

	var (
		user *user
		err  error
	)

	s := &struct {
		Password     string `json:"password"`
		RecoveryCode string `json:"code"`
	}{}

	if err = r.decodeJSON(s); err != nil {
		return err
	}

	if s.RecoveryCode == "" {
		if user, err = a.authenticate(r, true); err != nil {
			return err
		}
	} else {
		if user, err = a.ds.users.getByRecoveryCode(s.RecoveryCode); err != nil {
			return err
		}
		user.resetRecoveryCode()
	}

	if err = user.changePassword(s.Password); err != nil {
		return err
	}

	if err = a.ds.users.update(user); err != nil {
		return err
	}

	return renderString(w, http.StatusOK, "Password changed")
}

func (a *appContext) recoverPassword(w http.ResponseWriter, r *request) error {

	s := &struct {
		Email string `json:"email"`
	}{}

	if err := r.decodeJSON(s); err != nil {
		return err
	}
	if s.Email == "" {
		return &httpError{http.StatusBadRequest, "Missing email address"}
	}
	user, err := a.ds.users.getByEmail(s.Email)
	if err != nil {
		if isErrSqlNoRows(err) {
			return &httpError{http.StatusBadRequest, "Email address not found"}
		}
		return err
	}
	code, err := user.generateRecoveryCode()

	if err != nil {
		return err
	}

	if err := a.ds.users.update(user); err != nil {
		return err
	}

	go func() {
		if err := a.mailer.sendResetPasswordMail(user, code, r); err != nil {
			logError(err)
		}
	}()

	return renderString(w, http.StatusOK, "Password reset")
}
