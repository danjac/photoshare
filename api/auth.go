package api

import (
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"strings"
)

// lazily looks up user in session and stores in context.
var getCurrentUser = func(c web.C, r *http.Request) (*User, error) {

	obj, ok := c.Env["user"]
	if ok {
		return obj.(*User), nil
	}

	user, err := sessionMgr.GetCurrentUser(r)
	if err != nil {
		return nil, err
	}

	c.Env["user"] = user
	return user, nil
}

func checkAuth(c web.C, w http.ResponseWriter, r *http.Request) (*User, bool) {

	user, err := getCurrentUser(c, r)
	if err != nil {
		render.ServerError(w, err)
		return nil, false
	}
	if !user.IsAuthenticated {
		render.Status(w, http.StatusUnauthorized)
		return user, false
	}
	return user, true
}

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

func logout(c web.C, w http.ResponseWriter, r *http.Request) {

	user, ok := checkAuth(c, w, r)
	if !ok {
		return
	}

	if _, err := sessionMgr.Logout(w); err != nil {
		render.ServerError(w, err)
		return
	}

	sendMessage(&SocketMessage{user.Name, "", 0, "logout"})
	render.JSON(w, newSessionInfo(&User{}), http.StatusOK)

}

func authenticate(c web.C, w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(c, r)
	if err != nil {
		render.ServerError(w, err)
		return
	}

	render.JSON(w, newSessionInfo(user), http.StatusOK)
}

func login(c web.C, w http.ResponseWriter, r *http.Request) {

	s := &struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{}

	if err := parseJSON(r, s); err != nil {
		render.Error(w, http.StatusBadRequest)
		return
	}

	if s.Identifier == "" || s.Password == "" {
		render.String(w, "Missing login details", http.StatusBadRequest)
		return
	}

	user, err := userMgr.Authenticate(s.Identifier, s.Password)

	if err != nil {
		render.ServerError(w, err)
		return
	}
	if user == nil {
		render.String(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	if _, err := sessionMgr.Login(w, user); err != nil {
		render.ServerError(w, err)
		return
	}

	user.IsAuthenticated = true

	sendMessage(&SocketMessage{user.Name, "", 0, "login"})
	render.JSON(w, newSessionInfo(user), http.StatusOK)
}

func signup(c web.C, w http.ResponseWriter, r *http.Request) {

	s := &struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := parseJSON(r, s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &User{
		Name:     s.Name,
		Email:    strings.ToLower(s.Email),
		Password: s.Password,
	}

	validator := getUserValidator(user)

	if result, err := validator.Validate(); err != nil || !result.OK {
		if err != nil {
			render.ServerError(w, err)
			return
		}
		render.JSON(w, result, http.StatusBadRequest)
		return
	}

	if err := userMgr.Insert(user); err != nil {
		render.ServerError(w, err)
		return
	}

	if _, err := sessionMgr.Login(w, user); err != nil {
		render.ServerError(w, err)
		return
	}

	user.IsAuthenticated = true

	go func() {
		if err := sendWelcomeMail(user); err != nil {
			log.Println(err)
		}
	}()

	render.JSON(w, newSessionInfo(user), http.StatusOK)

}

func sendWelcomeMail(user *User) error {
	msg, err := MessageFromTemplate(
		"Welcome to photoshare!",
		[]string{user.Email},
		Config.Smtp.DefaultSender,
		signupTmpl,
		user,
	)
	if err != nil {
		return err
	}
	return mailer.Mail(msg)
}

func changePassword(c web.C, w http.ResponseWriter, r *http.Request) {

	var (
		user *User
		err  error
	)

	s := &struct {
		Password     string `json:"password"`
		RecoveryCode string `json:"code"`
	}{}

	if err = parseJSON(r, s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if s.RecoveryCode == "" {
		if user, err = getCurrentUser(c, r); err != nil {
			render.ServerError(w, err)
			return
		}
		if !user.IsAuthenticated {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		if user, err = userMgr.GetByRecoveryCode(s.RecoveryCode); err != nil {
			render.ServerError(w, err)
			return
		}
		if !user.IsAuthenticated {
			render.String(w, "Invalid code, no user found", http.StatusBadRequest)
			return
		}
		user.ResetRecoveryCode()
	}

	if err = user.ChangePassword(s.Password); err != nil {
		render.ServerError(w, err)
		return
	}

	if err = userMgr.Update(user); err != nil {
		render.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)

}

func recoverPassword(c web.C, w http.ResponseWriter, r *http.Request) {

	s := &struct {
		Email string `json:"email"`
	}{}

	if err := parseJSON(r, s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if s.Email == "" {
		render.String(w, "No email address provided", http.StatusBadRequest)
		return
	}
	user, err := userMgr.GetByEmail(s.Email)
	if err != nil {
		render.ServerError(w, err)
		return
	}
	if !user.IsAuthenticated {
		render.String(w, "No user found for this email address", http.StatusBadRequest)
		return
	}

	code, err := user.GenerateRecoveryCode()

	if err != nil {
		render.ServerError(w, err)
		return
	}

	if err := userMgr.Update(user); err != nil {
		render.ServerError(w, err)
		return
	}

	go func() {
		if err := sendResetPasswordMail(user, code, r); err != nil {
			log.Println(err)
		}
	}()

	render.Status(w, http.StatusOK)
}

func sendResetPasswordMail(user *User, recoveryCode string, r *http.Request) error {
	msg, err := MessageFromTemplate(
		"Reset your password",
		[]string{user.Email},
		Config.Smtp.DefaultSender,
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
	return mailer.Mail(msg)

}
