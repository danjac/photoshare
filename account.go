package photoshare

import (
	"net/http"
	"strings"
	"time"
)

func getAuthRedirectURL(ctx *context, w http.ResponseWriter, r *http.Request) error {

	url, err := ctx.auth.getRedirectURL(r, ctx.params.get("provider"))
	if err != nil {
		return err
	}
	return renderString(w, http.StatusOK, url)
}

func authCallback(ctx *context, w http.ResponseWriter, r *http.Request) error {

	info, err := ctx.auth.getUserInfo(r, ctx.params.get("provider"))
	if err != nil {
		return err
	}

	// tbd: handle new users
	user, err := ctx.datamapper.getUserByEmail(info.email)
	if err != nil {
		return err
	}

	authToken, err := ctx.session.createToken(user.ID)

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

func logout(ctx *context, w http.ResponseWriter, r *http.Request) error {

	if err := ctx.session.writeToken(w, 0); err != nil {
		return err
	}

	sendMessage(&socketMessage{ctx.user.Name, "", 0, "logout"})
	return renderJSON(w, newSessionInfo(&user{}), http.StatusOK)

}

func getSessionInfo(ctx *context, w http.ResponseWriter, r *http.Request) error {
	return renderJSON(w, newSessionInfo(ctx.user), http.StatusOK)
}

func login(ctx *context, w http.ResponseWriter, r *http.Request) error {

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

	user, err := ctx.datamapper.getUserByNameOrEmail(s.Identifier)
	if err != nil {
		if isErrSqlNoRows(err) {
			return invalidLogin
		}
		return err
	}
	if !user.checkPassword(s.Password) {
		return invalidLogin
	}

	if err := ctx.session.writeToken(w, user.ID); err != nil {
		return err
	}

	user.IsAuthenticated = true

	sendMessage(&socketMessage{user.Name, "", 0, "login"})
	return renderJSON(w, newSessionInfo(user), http.StatusCreated)
}

func signup(ctx *context, w http.ResponseWriter, r *http.Request) error {

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

	if err := ctx.validate(user, r); err != nil {
		return err
	}

	if err := ctx.datamapper.createUser(user); err != nil {
		return err
	}
	if err := ctx.session.writeToken(w, user.ID); err != nil {
		return err
	}

	user.IsAuthenticated = true

	go func() {
		if err := ctx.mailer.sendWelcomeMail(user); err != nil {
			logError(err)
		}
	}()

	return renderJSON(w, newSessionInfo(user), http.StatusCreated)

}

func changePassword(ctx *context, w http.ResponseWriter, r *http.Request) error {

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
		if user, err = ctx.authenticate(r, authLevelLogin); err != nil {
			return err
		}
	} else {
		if user, err = ctx.datamapper.getUserByRecoveryCode(s.RecoveryCode); err != nil {
			return err
		}
		user.resetRecoveryCode()
	}

	if err = user.changePassword(s.Password); err != nil {
		return err
	}
	if err := ctx.validate(user, r); err != nil {
		return err
	}
	if err := ctx.datamapper.updateUser(user); err != nil {
		return err
	}

	return renderString(w, http.StatusOK, "Password changed")
}

func recoverPassword(ctx *context, w http.ResponseWriter, r *http.Request) error {

	s := &struct {
		Email string `json:"email"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		return err
	}
	if s.Email == "" {
		return httpError{http.StatusBadRequest, "Missing email address"}
	}
	user, err := ctx.datamapper.getUserByEmail(s.Email)
	if err != nil {
		if isErrSqlNoRows(err) {
			return httpError{http.StatusBadRequest, "Email address not found"}
		}
		return err
	}
	code, err := user.generateRecoveryCode()

	if err := ctx.datamapper.updateUser(user); err != nil {
		return err
	}

	go func() {
		if err := ctx.mailer.sendResetPasswordMail(user, code, r); err != nil {
			logError(err)
		}
	}()

	return renderString(w, http.StatusOK, "Password reset")
}
