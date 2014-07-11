package api

import (
	"log"
	"net/http"
	"strings"
)

var getCurrentUser = func(r *http.Request) (*User, error) {

	user, err := sessionMgr.GetCurrentUser(r)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// gets current user. If user not authenticated, writes a 401 error. Returns true
// if no error/user authenticated.
func getUserOr401(w http.ResponseWriter, r *http.Request) (*User, bool) {

	user, err := getCurrentUser(r)
	if err != nil {
		serverError(w, err)
		return nil, false
	}
	if !user.IsAuthenticated {
		http.Error(w, "You must be logged in", http.StatusUnauthorized)
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

func logout(w http.ResponseWriter, r *http.Request) {

	user, ok := getUserOr401(w, r)
	if !ok {
		return
	}

	if _, err := sessionMgr.Logout(w); err != nil {
		serverError(w, err)
		return
	}

	sendMessage(&SocketMessage{user.Name, "", 0, "logout"})
	writeJSON(w, newSessionInfo(&User{}), http.StatusOK)

}

func authenticate(w http.ResponseWriter, r *http.Request) {

	user, err := getCurrentUser(r)
	if err != nil {
		serverError(w, err)
		return
	}

	writeJSON(w, newSessionInfo(user), http.StatusOK)
}

func login(w http.ResponseWriter, r *http.Request) {

	s := &struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		http.Error(w, "Invalid login details", http.StatusBadRequest)
		return
	}

	if s.Identifier == "" || s.Password == "" {
		http.Error(w, "Missing login details", http.StatusBadRequest)
		return
	}

	user, exists, err := userMgr.Authenticate(s.Identifier, s.Password)

	if err != nil {
		serverError(w, err)
		return
	}
	if !exists {
		http.Error(w, "Invalid email or password", http.StatusBadRequest)
		return
	}

	if _, err := sessionMgr.Login(w, user); err != nil {
		serverError(w, err)
		return
	}

	user.IsAuthenticated = true

	sendMessage(&SocketMessage{user.Name, "", 0, "login"})
	writeJSON(w, newSessionInfo(user), http.StatusCreated)
}

func signup(w http.ResponseWriter, r *http.Request) {

	s := &struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		http.Error(w, "Invalid signup data", http.StatusBadRequest)
		return
	}

	user := &User{
		Name:     s.Name,
		Email:    strings.ToLower(s.Email),
		Password: s.Password,
	}

	validator := getUserValidator(user)

	if result, err := formHandler.Validate(validator); err != nil || !result.OK {
		if err != nil {
			serverError(w, err)
			return
		}
		result.Write(w)
		return
	}

	if err := userMgr.Insert(user); err != nil {
		serverError(w, err)
		return
	}

	if _, err := sessionMgr.Login(w, user); err != nil {
		serverError(w, err)
		return
	}

	user.IsAuthenticated = true

	go func() {
		if err := sendWelcomeMail(user); err != nil {
			log.Println(err)
		}
	}()

	writeJSON(w, newSessionInfo(user), http.StatusCreated)

}

func sendWelcomeMail(user *User) error {
	msg, err := MessageFromTemplate(
		"Welcome to photoshare!",
		[]string{user.Email},
		config.SmtpDefaultSender,
		signupTmpl,
		user,
	)
	if err != nil {
		return err
	}
	return mailer.Mail(msg)
}

func changePassword(w http.ResponseWriter, r *http.Request) {

	var (
		user *User
		err  error
		ok   bool
	)

	s := &struct {
		Password     string `json:"password"`
		RecoveryCode string `json:"code"`
	}{}

	if err = decodeJSON(r, s); err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	if s.RecoveryCode == "" {
		if user, ok = getUserOr401(w, r); !ok {
			return
		}
	} else {
		if user, ok, err = userMgr.GetByRecoveryCode(s.RecoveryCode); err != nil {
			serverError(w, err)
			return
		}
		if !ok {
			http.Error(w, "Invalid code, no user found", http.StatusBadRequest)
			return
		}
		user.ResetRecoveryCode()
	}

	if err = user.ChangePassword(s.Password); err != nil {
		serverError(w, err)
		return
	}

	if err = userMgr.Update(user); err != nil {
		serverError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func recoverPassword(w http.ResponseWriter, r *http.Request) {

	s := &struct {
		Email string `json:"email"`
	}{}

	if err := decodeJSON(r, s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if s.Email == "" {
		http.Error(w, "No email address provided", http.StatusBadRequest)
		return
	}
	user, exists, err := userMgr.GetByEmail(s.Email)
	if err != nil {
		serverError(w, err)
		return
	}
	if !exists {
		http.Error(w, "No user found for this email address", http.StatusBadRequest)
		return
	}

	code, err := user.GenerateRecoveryCode()

	if err != nil {
		serverError(w, err)
		return
	}

	if err := userMgr.Update(user); err != nil {
		serverError(w, err)
		return
	}

	go func() {
		if err := sendResetPasswordMail(user, code, r); err != nil {
			log.Println(err)
		}
	}()

	w.WriteHeader(http.StatusNoContent)
}

func sendResetPasswordMail(user *User, recoveryCode string, r *http.Request) error {
	msg, err := MessageFromTemplate(
		"Reset your password",
		[]string{user.Email},
		config.SmtpDefaultSender,
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
