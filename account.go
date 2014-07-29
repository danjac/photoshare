package photoshare

import (
	"net/http"
	"strings"
	"time"
)

func getAuthRedirectURL(c *context, w http.ResponseWriter, r *http.Request) error {

	url, err := c.auth.getRedirectURL(r, c.params.get("provider"))
	if err != nil {
		return err
	}
	return renderString(w, http.StatusOK, url)
}

func authCallback(c *context, w http.ResponseWriter, r *http.Request) error {

	info, err := c.auth.getUserInfo(r, c.params.get("provider"))
	if err != nil {
		return err
	}

	// tbd: handle new users
	user, err := c.ds.getUserByEmail(info.email)
	if err != nil {
		return err
	}

	authToken, err := c.session.createToken(user.ID)

	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:    "authToken",
		Value:   authToken,
		Path:    "/",
		Expires: time.Now().AddDate(0, 0, 1),
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return nil
}

func logout(c *context, w http.ResponseWriter, r *http.Request) error {

	u, err := c.getUser(r, true)
	if err != nil {
		return err
	}
	if err := c.session.writeToken(w, 0); err != nil {
		return err
	}

	sendMessage(&socketMessage{u.Name, "", 0, "logout"})
	return renderJSON(w, newSessionInfo(&user{}), http.StatusOK)

}

func getSessionInfo(c *context, w http.ResponseWriter, r *http.Request) error {

	user, err := c.getUser(r, false)
	if err != nil {
		return err
	}
	return renderJSON(w, newSessionInfo(user), http.StatusOK)
}

func login(c *context, w http.ResponseWriter, r *http.Request) error {

	s := &struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{}

	var invalidLogin = httpError{http.StatusBadRequest, "Invalid email or password"}

	if err := decodeJSON(r, s); err != nil {
		return err
	}

	if s.Identifier == "" || s.Password == "" {
		return invalidLogin
	}

	user, err := c.ds.getUserByNameOrEmail(s.Identifier)
	if err != nil {
		if isErrSqlNoRows(err) {
			return invalidLogin
		}
		return err
	}
	if !user.checkPassword(s.Password) {
		return invalidLogin
	}

	if err := c.session.writeToken(w, user.ID); err != nil {
		return err
	}

	user.IsAuthenticated = true

	sendMessage(&socketMessage{user.Name, "", 0, "login"})
	return renderJSON(w, newSessionInfo(user), http.StatusCreated)
}

func signup(c *context, w http.ResponseWriter, r *http.Request) error {

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

	if err := c.validate(user); err != nil {
		return err
	}

	if err := c.ds.createUser(user); err != nil {
		return err
	}
	if err := c.session.writeToken(w, user.ID); err != nil {
		return err
	}

	user.IsAuthenticated = true

	go func() {
		if err := c.mailer.sendWelcomeMail(user); err != nil {
			logError(err)
		}
	}()

	return renderJSON(w, newSessionInfo(user), http.StatusCreated)

}

func changePassword(c *context, w http.ResponseWriter, r *http.Request) error {

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
		if user, err = c.getUser(r, true); err != nil {
			return err
		}
	} else {
		if user, err = c.ds.getUserByRecoveryCode(s.RecoveryCode); err != nil {
			return err
		}
		user.resetRecoveryCode()
	}

	if err = user.changePassword(s.Password); err != nil {
		return err
	}
	if err := c.validate(user); err != nil {
		return err
	}
	if err := c.ds.updateUser(user); err != nil {
		return err
	}

	return renderString(w, http.StatusOK, "Password changed")
}

func recoverPassword(c *context, w http.ResponseWriter, r *http.Request) error {

	s := &struct {
		Email string `json:"email"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		return err
	}
	if s.Email == "" {
		return httpError{http.StatusBadRequest, "Missing email address"}
	}
	user, err := c.ds.getUserByEmail(s.Email)
	if err != nil {
		if isErrSqlNoRows(err) {
			return httpError{http.StatusBadRequest, "Email address not found"}
		}
		return err
	}
	code, err := user.generateRecoveryCode()

	if err := c.ds.updateUser(user); err != nil {
		return err
	}

	go func() {
		if err := c.mailer.sendResetPasswordMail(user, code, r); err != nil {
			logError(err)
		}
	}()

	return renderString(w, http.StatusOK, "Password reset")
}
