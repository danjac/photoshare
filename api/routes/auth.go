package routes

import (
	"github.com/danjac/photoshare/api/config"
	"github.com/danjac/photoshare/api/email"
	"github.com/danjac/photoshare/api/models"
	"github.com/danjac/photoshare/api/session"
	"github.com/danjac/photoshare/api/validation"
	"github.com/zenazn/goji/web"
	"log"
	"net/http"
	"strings"
	"text/template"
)

var sessionMgr = session.NewSessionManager()

var (
	signupTmpl,
	recoverPassTmpl *template.Template
)

var getUserValidator = func(user *models.User) validation.Validator {
	return validation.NewUserValidator(user)
}

// lazily looks up user in session and stores in context.
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
type sessionInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsAdmin  bool   `json:"isAdmin"`
	LoggedIn bool   `json:"loggedIn"`
}

func newSessionInfo(user *models.User) *sessionInfo {
	if user == nil || user.ID == 0 || !user.IsAuthenticated {
		return &sessionInfo{}
	}

	return &sessionInfo{user.ID, user.Name, user.IsAdmin, true}
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
		w.WriteHeader(http.StatusBadRequest)
		return
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

	s := &struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	if err := parseJSON(r, s); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &models.User{
		Name:     s.Name,
		Email:    strings.ToLower(s.Email),
		Password: s.Password,
	}

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

	if msg, err := email.MessageFromTemplate(
		"Welcome to photoshare!",
		[]string{user.Email},
		config.DefaultEmailSender,
		signupTmpl,
		user,
	); err == nil {
		go func() {
			if err := mailer.Mail(msg); err != nil {
				log.Println(err)
			}
		}()

	} else {
		panic(err)
	}

	writeJSON(w, newSessionInfo(user), http.StatusOK)

}

func changePassword(c web.C, w http.ResponseWriter, r *http.Request) {

	var (
		user *models.User
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
			panic(err)
		}
		if !user.IsAuthenticated {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else {
		if user, err = userMgr.GetByRecoveryCode(s.RecoveryCode); err != nil {
			panic(err)
		}
		if !user.IsAuthenticated {
			writeString(w, "No user found for this email address", http.StatusBadRequest)
			return
		}
		user.ResetRecoveryCode()
	}

	if err = user.ChangePassword(s.Password); err != nil {
		panic(err)
	}

	if err = userMgr.Update(user); err != nil {
		panic(err)
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
		writeString(w, "No email address provided", http.StatusBadRequest)
		return
	}
	user, err := userMgr.GetByEmail(s.Email)
	if err != nil {
		panic(err)
	}
	if !user.IsAuthenticated {
		http.NotFound(w, r)
		return
	}

	code, err := user.GenerateRecoveryCode()
	if err != nil {
		panic(err)
	}

	if err := userMgr.Update(user); err != nil {
		panic(err)
	}
	if msg, err := email.MessageFromTemplate(
		"Reset your password",
		[]string{user.Email},
		config.DefaultEmailSender,
		recoverPassTmpl,
		&struct {
			Name         string
			RecoveryCode string
			Url          string
		}{
			user.Name,
			code,
			baseURL(r),
		},
	); err == nil {
		go func() {
			if err := mailer.Mail(msg); err != nil {
				log.Println(err)
			}
		}()
	} else {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func init() {
	signupTmpl = parseTemplate("signup.tmpl")
	recoverPassTmpl = parseTemplate("recover_pass.tmpl")
}
