package api

import (
	"github.com/zenazn/goji/web"
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

func (a *appContext) authenticate(c web.C, r *http.Request, required bool) (*user, error) {

	var (
		u            *user
		invalidLogin error
	)

	if required {
		invalidLogin = httpError(http.StatusUnauthorized, "You must be logged in")
	}

	obj, ok := c.Env["user"]

	if ok {
		u = obj.(*user)
	} else {
		userID, err := a.sessionMgr.readToken(r)
		if err != nil {
			return u, err
		}
		if userID == 0 {
			return u, invalidLogin
		}
		u, err = a.ds.users.getActive(userID)
		if err != nil {
			if isErrSqlNoRows(err) {
				return u, invalidLogin
			}
			return u, err
		}
		c.Env["user"] = u
	}
	u.IsAuthenticated = true

	return u, nil
}

func (a *appContext) logout(c web.C, w http.ResponseWriter, r *http.Request) error {

	u, err := a.authenticate(c, r, true)
	if err != nil {
		return err
	}

	if err := a.sessionMgr.writeToken(w, 0); err != nil {
		return err
	}

	sendMessage(&SocketMessage{u.Name, "", 0, "logout"})
	return renderJSON(w, newSessionInfo(&user{}), http.StatusOK)

}

func (a *appContext) getSessionInfo(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, false)
	if err != nil {
		return err
	}

	return renderJSON(w, newSessionInfo(user), http.StatusOK)
}

func (a *appContext) login(_ web.C, w http.ResponseWriter, r *http.Request) error {

	s := &struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{}

	var invalidLogin = httpError(http.StatusBadRequest, "Invalid email or password")

	if err := decodeJSON(r, s); err != nil {
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

	sendMessage(&SocketMessage{user.Name, "", 0, "login"})
	return renderJSON(w, newSessionInfo(user), http.StatusCreated)
}

func (a *appContext) signup(c web.C, w http.ResponseWriter, r *http.Request) error {

	s := &struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := decodeJSON(r, s); err != nil {
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

	if err := a.ds.users.insert(user); err != nil {
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

func (a *appContext) changePassword(c web.C, w http.ResponseWriter, r *http.Request) error {

	var (
		user *user
		err  error
	)

	s := &struct {
		Password     string `json:"password"`
		RecoveryCode string `json:"code"`
	}{}

	if err = decodeJSON(r, s); err != nil {
		return err
	}

	if s.RecoveryCode == "" {
		if user, err = a.authenticate(c, r, true); err != nil {
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

func (a *appContext) recoverPassword(_ web.C, w http.ResponseWriter, r *http.Request) error {

	s := &struct {
		Email string `json:"email"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		return err
	}
	if s.Email == "" {
		return httpError(http.StatusBadRequest, "Missing email address")
	}
	user, err := a.ds.users.getByEmail(s.Email)
	if err != nil {
		if isErrSqlNoRows(err) {
			return httpError(http.StatusBadRequest, "Email address not found")
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
