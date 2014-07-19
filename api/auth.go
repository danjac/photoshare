package api

import (
	"database/sql"
	"github.com/zenazn/goji/web"
	"log"
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

func newSessionInfo(user *User) *sessionInfo {
	if user == nil || user.ID == 0 || !user.IsAuthenticated {
		return &sessionInfo{}
	}

	return &sessionInfo{user.ID, user.Name, user.IsAdmin, true}
}

func (a *AppContext) authenticate(c web.C, r *http.Request, required bool) (*User, error) {

	var user *User
	var invalidLogin = httpError(http.StatusUnauthorized, "You must be logged in")

	obj, ok := c.Env["user"]

	if ok {
		user = obj.(*User)
	} else {
		userID, err := a.sessionMgr.ReadToken(r)
		if err != nil {
			return user, err
		}
		if userID == 0 {
			return user, invalidLogin
		}
		user, err = a.userDS.GetActive(userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return user, invalidLogin
			}
			return user, err
		}
		c.Env["user"] = user
	}
	user.IsAuthenticated = true

	return user, nil
}

func (a *AppContext) logout(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, true)
	if err != nil {
		return err
	}

	if err := a.sessionMgr.WriteToken(w, 0); err != nil {
		return err
	}

	sendMessage(&SocketMessage{user.Name, "", 0, "logout"})
	return renderJSON(w, newSessionInfo(&User{}), http.StatusOK)

}

func (a *AppContext) getSessionInfo(c web.C, w http.ResponseWriter, r *http.Request) error {

	user, err := a.authenticate(c, r, false)
	if err != nil {
		return err
	}

	return renderJSON(w, newSessionInfo(user), http.StatusOK)
}

func (a *AppContext) login(_ web.C, w http.ResponseWriter, r *http.Request) error {

	s := &struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		return httpError(http.StatusBadRequest, "Invalid data")
	}

	if s.Identifier == "" || s.Password == "" {
		return httpError(http.StatusBadRequest, "Missing email or password")
	}

	user, err := a.userDS.Authenticate(s.Identifier, s.Password)
	if err != nil {
		return err
	}

	if err := a.sessionMgr.WriteToken(w, user.ID); err != nil {
		return err
	}

	user.IsAuthenticated = true

	sendMessage(&SocketMessage{user.Name, "", 0, "login"})
	return renderJSON(w, newSessionInfo(user), http.StatusCreated)
}

func (a *AppContext) signup(c web.C, w http.ResponseWriter, r *http.Request) error {

	s := &struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		return err
	}

	user := &User{
		Name:     s.Name,
		Email:    strings.ToLower(s.Email),
		Password: s.Password,
	}

	validator := NewUserValidator(user, a.userDS)

	if err := validate(validator); err != nil {
		return err
	}

	if err := a.userDS.Insert(user); err != nil {
		return err
	}

	if err := a.sessionMgr.WriteToken(w, user.ID); err != nil {
		return err
	}

	user.IsAuthenticated = true

	go func() {
		if err := a.sendWelcomeMail(user); err != nil {
			log.Println(err)
		}
	}()

	return renderJSON(w, newSessionInfo(user), http.StatusCreated)

}

func (a *AppContext) changePassword(c web.C, w http.ResponseWriter, r *http.Request) error {

	var (
		user *User
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
		if user, err = a.userDS.GetByRecoveryCode(s.RecoveryCode); err != nil {
			return err
		}
		user.ResetRecoveryCode()
	}

	if err = user.ChangePassword(s.Password); err != nil {
		return err
	}

	if err = a.userDS.Update(user); err != nil {
		return err
	}

	return renderString(w, http.StatusOK, "Password changed")
}

func (a *AppContext) recoverPassword(_ web.C, w http.ResponseWriter, r *http.Request) error {

	s := &struct {
		Email string `json:"email"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		return err
	}
	if s.Email == "" {
		return httpError(http.StatusBadRequest, "Missing email address")
	}
	user, err := a.userDS.GetByEmail(s.Email)
	if err != nil {
		return err
	}
	code, err := user.GenerateRecoveryCode()

	if err != nil {
		return err
	}

	if err := a.userDS.Update(user); err != nil {
		return err
	}

	go func() {
		if err := a.sendResetPasswordMail(user, code, r); err != nil {
			log.Println(err)
		}
	}()

	return renderString(w, http.StatusOK, "Password reset")
}

func (a *AppContext) sendResetPasswordMail(user *User, recoveryCode string, r *http.Request) error {
	msg, err := MessageFromTemplate(
		"Reset your password",
		[]string{user.Email},
		a.config.SmtpDefaultSender,
		recoverPassTmpl,
		&struct {
			Name         string
			RecoveryCode string
			Url          string
		}{
			user.Name,
			recoveryCode,
			baseURL(r),
		},
	)
	if err != nil {
		return err
	}
	return a.mailer.Mail(msg)
}

func (a *AppContext) sendWelcomeMail(user *User) error {
	msg, err := MessageFromTemplate(
		"Welcome to photoshare!",
		[]string{user.Email},
		a.config.SmtpDefaultSender,
		signupTmpl,
		user,
	)
	if err != nil {
		return err
	}
	return a.mailer.Mail(msg)
}
